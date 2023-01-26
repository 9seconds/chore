package cli

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"text/tabwriter"
	"unicode"

	"github.com/9seconds/chore/internal/script"
	"github.com/spf13/cobra"
)

const (
	DirSizeUnknown = "?"
	DirSizeError   = "ERROR"

	ByteBase = 1024

	RequiredTrue  = "✔"
	RequiredFalse = "✖"

	TabSize = 8
)

var byteUnits = [6]string{
	"B",
	"KB",
	"MB",
	"GB",
	"TB",
	"EB",
}

func NewShow() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "show [flags] [namespace] [script]",
		Aliases: []string{"s"},
		Short:   "Show details on scripts and namespaces.",
		Args: cobra.MatchAll(
			cobra.MaximumNArgs(2), //nolint: gomnd
			argumentOptional(0, validAsciiName(0, ErrNamespaceInvalid)),
			argumentOptional(1, validAsciiName(1, ErrScriptInvalid)),
			argumentOptional(0, validNamespace(0)),
			argumentOptional(1, validScript(0, 1)),
		),
		RunE:              mainShow,
		ValidArgsFunction: completeNamespaceScript,
	}

	flags := cmd.Flags()

	flags.BoolP("show-path", "p", false, "show path to the script")
	flags.BoolP("show-config", "c", false, "show path to the script config")
	flags.BoolP("show-data", "t", false, "show path to the script data directory")
	flags.BoolP("show-cache", "a", false, "show path to the script config directory")
	flags.BoolP("show-state", "s", false, "show path to the script state directory")
	flags.BoolP("show-runtime", "r", false, "show path to the script runtime directory")

	return cmd
}

func mainShow(cmd *cobra.Command, args []string) error { //nolint: cyclop
	switch len(args) {
	case 0:
		return mainShowListNamespaces(cmd)
	case 1:
		return mainShowListScripts(cmd, args[0])
	}

	scr := &script.Script{
		Namespace:  args[0],
		Executable: args[1],
	}

	if err := scr.Init(); err != nil {
		return fmt.Errorf("cannot initialize script: %w", err)
	}

	defer scr.Cleanup()

	showPaths, err := cmd.Flags().GetBool("show-path")
	if err != nil {
		return err
	}

	showConfig, err := cmd.Flags().GetBool("show-config")
	if err != nil {
		return err
	}

	showData, err := cmd.Flags().GetBool("show-data")
	if err != nil {
		return err
	}

	showCache, err := cmd.Flags().GetBool("show-cache")
	if err != nil {
		return err
	}

	showState, err := cmd.Flags().GetBool("show-state")
	if err != nil {
		return err
	}

	showRuntime, err := cmd.Flags().GetBool("show-runtime")
	if err != nil {
		return err
	}

	if showPaths || showConfig || showData || showCache || showState || showRuntime {
		mainShowFlags(cmd, scr, showPaths, showConfig, showData, showCache, showState, showRuntime)
	} else {
		mainShowTables(cmd, scr)
	}

	return nil
}

func mainShowListNamespaces(cmd *cobra.Command) error {
	names, err := script.ListNamespaces()
	if err != nil {
		return fmt.Errorf("cannot list namespaces: %w", err)
	}

	for _, name := range names {
		cmd.Println(name)
	}

	return nil
}

func mainShowListScripts(cmd *cobra.Command, namespace string) error {
	names, err := script.ListScripts(namespace)
	if err != nil {
		return fmt.Errorf("cannot list scripts: %w", err)
	}

	for _, name := range names {
		cmd.Println(name)
	}

	return nil
}

func mainShowFlags(cmd *cobra.Command, scr *script.Script, paths, config, data, cache, state, runtime bool) {
	if paths {
		cmd.Println(scr.Path())
	}

	if config {
		cmd.Println(scr.ConfigPath())
	}

	if data {
		cmd.Println(scr.DataPath())
	}

	if cache {
		cmd.Println(scr.CachePath())
	}

	if state {
		cmd.Println(scr.StatePath())
	}

	if runtime {
		cmd.Println(scr.RuntimePath())
	}
}

func mainShowTables(cmd *cobra.Command, scr *script.Script) {
	ctx := cmd.Context()
	dirSizes := map[string]*atomic.Int64{
		scr.DataPath():    {},
		scr.StatePath():   {},
		scr.CachePath():   {},
		scr.RuntimePath(): {},
	}

	waiters := &sync.WaitGroup{}

	waiters.Add(len(dirSizes))

	for path, accumulator := range dirSizes {
		go func(seedPath string, accumulator *atomic.Int64) {
			defer waiters.Done()

			err := filepath.Walk(seedPath, func(_ string, info fs.FileInfo, _ error) error {
				accumulator.Add(info.Size())

				select {
				case <-ctx.Done():
					return fmt.Errorf("cancelled because of context: %w", ctx.Err())
				default:
					return nil
				}
			})
			if err != nil {
				log.Printf("cannot complete traversal of %s: %v", seedPath, err)

				accumulator.Store(-1)
			}
		}(path, accumulator)
	}

	waiters.Wait()

	buf := &strings.Builder{}

	mainShowDescription(buf, scr)
	mainShowMainTable(buf, scr, dirSizes)
	mainShowTableParameters(buf, scr)
	mainShowTableFlags(buf, scr)

	cmd.Println(strings.TrimRightFunc(buf.String(), unicode.IsSpace))
}

func mainShowDescription(buf io.Writer, scr *script.Script) {
	conf := scr.Config()

	if conf.Description != "" {
		io.WriteString(buf, strings.TrimSpace(conf.Description)) //nolint: errcheck
		io.WriteString(buf, "\n\n")                              //nolint: errcheck
	}
}

func mainShowMainTable(buf io.Writer, scr *script.Script, dirSizes map[string]*atomic.Int64) {
	defer io.WriteString(buf, "\n") //nolint: errcheck

	conf := scr.Config()

	writer := mainTabwriter(buf)

	defer writer.Flush()

	fmt.Fprintf(writer, "Path:\t%s\n", scr.Path())
	fmt.Fprintf(writer, "Config path:\t%s\n", scr.ConfigPath())
	fmt.Fprintf(writer, "Data path:\t%s\t%s\n", scr.DataPath(), mainShowDirSize(dirSizes[scr.DataPath()]))
	fmt.Fprintf(writer, "Cache path:\t%s\t%s\n", scr.CachePath(), mainShowDirSize(dirSizes[scr.CachePath()]))
	fmt.Fprintf(writer, "State path:\t%s\t%s\n", scr.StatePath(), mainShowDirSize(dirSizes[scr.StatePath()]))
	fmt.Fprintf(writer, "Runtime path:\t%s\t%s\n\n", scr.RuntimePath(), mainShowDirSize(dirSizes[scr.RuntimePath()]))

	fmt.Fprintf(writer, "Network:\t%s\n", strconv.FormatBool(conf.Network))
	fmt.Fprintf(writer, "Git:\t%s\n", conf.Git.String())
}

func mainShowTableParameters(buf io.Writer, scr *script.Script) {
	conf := scr.Config()

	if len(conf.Parameters) == 0 {
		return
	}

	defer io.WriteString(buf, "\n") //nolint: errcheck

	names := make([]string, 0, len(conf.Parameters))

	for k := range conf.Parameters {
		names = append(names, k)
	}

	// sort by (required_as_int, name_as_str)
	sort.Slice(names, func(i, j int) bool {
		nameI := names[i]
		nameJ := names[j]
		paramI := conf.Parameters[nameI]
		paramJ := conf.Parameters[nameJ]
		valueI := mainShowBoolToInt(paramI.Required())
		valueJ := mainShowBoolToInt(paramJ.Required())

		if valueI == valueJ {
			return nameI < nameJ
		}

		return valueI > valueJ
	})

	writer := mainTabwriter(buf)

	defer writer.Flush()

	fmt.Fprintln(writer, "Parameter\tDescription\tRequired?\tType\tSpecification")
	fmt.Fprintln(writer, "╴╴╴╴╴╴╴╴╴\t╴╴╴╴╴╴╴╴╴╴╴\t╴╴╴╴╴╴╴╴╴\t╴╴╴╴\t╴╴╴╴╴╴╴╴╴╴╴╴╴")

	for _, name := range names {
		param := conf.Parameters[name]

		fmt.Fprintf(
			writer,
			"%s\t%s\t%s\t%s\t%s\n",
			name,
			param.Description(),
			mainShowRequired(param.Required()),
			param.Type(),
			mainShowParameterSpec(param.Specification()))
	}
}

func mainShowTableFlags(buf io.Writer, scr *script.Script) {
	conf := scr.Config()

	if len(conf.Flags) == 0 {
		return
	}

	names := make([]string, 0, len(conf.Flags))

	for k := range conf.Flags {
		names = append(names, k)
	}

	// sort by (required_as_int, name_as_str)
	sort.Slice(names, func(i, j int) bool {
		nameI := names[i]
		nameJ := names[j]
		flagI := conf.Flags[nameI]
		flagJ := conf.Flags[nameJ]
		valueI := mainShowBoolToInt(flagI.Required())
		valueJ := mainShowBoolToInt(flagJ.Required())

		if valueI == valueJ {
			return nameI < nameJ
		}

		return valueI > valueJ
	})

	writer := mainTabwriter(buf)

	defer writer.Flush()

	fmt.Fprintln(writer, "Flag\tDescription\tRequired?")
	fmt.Fprintln(writer, "╴╴╴╴\t╴╴╴╴╴╴╴╴╴╴╴\t╴╴╴╴╴╴╴╴╴")

	for _, name := range names {
		flag := conf.Flags[name]

		fmt.Fprintf(
			writer,
			"%s\t%s\t%s\n",
			name,
			flag.Description(),
			mainShowRequired(flag.Required()))
	}
}

func mainShowDirSize(atomicValue *atomic.Int64) string {
	sizeInBytes := atomicValue.Load()

	if sizeInBytes < 0 {
		return DirSizeError
	}

	size := float64(sizeInBytes)
	unit := 0

	for size >= ByteBase && unit < len(byteUnits) {
		size /= ByteBase
		unit++
	}

	return strconv.FormatFloat(size, 'g', 2, 64) + byteUnits[unit]
}

func mainShowBoolToInt(value bool) int {
	if value {
		return 1
	}

	return 0
}

func mainShowRequired(value bool) string {
	if value {
		return RequiredTrue
	}

	return RequiredFalse
}

func mainShowParameterSpec(spec map[string]string) string {
	names := make([]string, 0, len(spec))

	for name := range spec {
		names = append(names, name)
	}

	sort.Strings(names)

	kvs := make([]string, 0, len(spec))

	for _, name := range names {
		kvs = append(kvs, fmt.Sprintf("%s=%q", name, spec[name]))
	}

	return strings.Join(kvs, " ")
}

func mainTabwriter(writer io.Writer) *tabwriter.Writer {
	return tabwriter.NewWriter(writer, 0, TabSize, 1, '\t', 0)
}
