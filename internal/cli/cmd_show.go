package cli

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/9seconds/chore/internal/script"
	"github.com/cheynewallace/tabby"
	"github.com/spf13/cobra"
)

const (
	DirSizeUnknown = "?"
	DirSizeError   = "ERROR"

	ByteBase = 1024

	RequiredTrue  = "✔"
	RequiredFalse = "✖"
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
			argumentOptional(0, validScriptName(0, ErrNamespaceInvalid)),
			argumentOptional(1, validScriptName(1, ErrScriptInvalid)),
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
		return mainShowListNamespaces()
	case 1:
		return mainShowListScripts(args[0])
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
		mainShowFlags(scr, showPaths, showConfig, showData, showCache, showState, showRuntime)
	} else {
		mainShowTables(cmd.Context(), scr)
	}

	return nil
}

func mainShowListNamespaces() error {
	names, err := script.ListNamespaces()
	if err != nil {
		return fmt.Errorf("cannot list namespaces: %w", err)
	}

	for _, name := range names {
		fmt.Println(name)
	}

	return nil
}

func mainShowListScripts(namespace string) error {
	names, err := script.ListScripts(namespace)
	if err != nil {
		return fmt.Errorf("cannot list scripts: %w", err)
	}

	for _, name := range names {
		fmt.Println(name)
	}

	return nil
}

func mainShowFlags(scr *script.Script, paths, config, data, cache, state, runtime bool) {
	if paths {
		fmt.Println(scr.Path())
	}

	if config {
		fmt.Println(scr.ConfigPath())
	}

	if data {
		fmt.Println(scr.DataPath())
	}

	if cache {
		fmt.Println(scr.CachePath())
	}

	if state {
		fmt.Println(scr.StatePath())
	}

	if runtime {
		fmt.Println(scr.RuntimePath())
	}
}

func mainShowTables(ctx context.Context, scr *script.Script) {
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

	mainShowTablePaths(scr, dirSizes)
	fmt.Println()

	mainShowTableGlobals(scr)
	mainShowTableParameters(scr)
	mainShowTableFlags(scr)
}

func mainShowTablePaths(scr *script.Script, dirSizes map[string]*atomic.Int64) {
	table := tabby.New()

	table.AddLine("Path:", scr.Path(), "")
	table.AddLine("Config path:", scr.ConfigPath(), "")
	table.AddLine("Data path:", scr.DataPath(), mainShowDirSize(dirSizes[scr.DataPath()]))
	table.AddLine("Cache path:", scr.CachePath(), mainShowDirSize(dirSizes[scr.CachePath()]))
	table.AddLine("State path:", scr.StatePath(), mainShowDirSize(dirSizes[scr.StatePath()]))
	table.AddLine("Runtime path:", scr.RuntimePath(), mainShowDirSize(dirSizes[scr.RuntimePath()]))
	table.Print()
}

func mainShowTableGlobals(scr *script.Script) {
	table := tabby.New()
	conf := scr.Config()

	table.AddLine("Network:", strconv.FormatBool(conf.Network))
	table.AddLine("Git:", conf.Git.String())

	table.Print()
}

func mainShowTableParameters(scr *script.Script) {
	table := tabby.New()
	conf := scr.Config()

	if len(conf.Parameters) == 0 {
		return
	}

	fmt.Println()

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

	table.AddHeader("Parameter", "Description", "Required?", "Type", "Specification")

	for _, name := range names {
		param := conf.Parameters[name]

		table.AddLine(
			name,
			param.Description(),
			mainShowRequired(param.Required()),
			param.Type(),
			mainShowParameterSpec(param.Specification()))
	}

	table.Print()
}

func mainShowTableFlags(scr *script.Script) {
	table := tabby.New()
	conf := scr.Config()

	if len(conf.Flags) == 0 {
		return
	}

	fmt.Println()

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

	table.AddHeader("Flag", "Description", "Required?")

	for _, name := range names {
		flag := conf.Flags[name]

		table.AddLine(
			name,
			flag.Description(),
			mainShowRequired(flag.Required()))
	}

	table.Print()
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

	return strconv.FormatFloat(size, 'f', 2, 64) + byteUnits[unit]
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
