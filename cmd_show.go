package main

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

	"github.com/9seconds/chore/internal/cli"
	"github.com/9seconds/chore/internal/script"
	"github.com/cheynewallace/tabby"
)

const (
	DirSizeUnknown = "?"
	DirSizeError   = "ERROR"

	ByteBase = 1024

	RequiredTrue  = "✔"
	RequiredFalse = "✖"
)

var ByteUnits = [6]string{
	"B",
	"KB",
	"MB",
	"GB",
	"TB",
	"EB",
}

type CliCmdShow struct {
	Namespace cli.Namespace `arg:"" optional:"" help:"Script namespace. Dot takes one from environment variable CHORE_NAMESPACE."`
	Script    string        `arg:"" optional:"" help:"Script name."`

	ShowPath    bool `short:"p" help:"Show path to the script."`
	ShowConfig  bool `short:"c" help:"Show path to the script config."`
	ShowData    bool `short:"t" help:"Show path to the script data directory."`
	ShowCache   bool `short:"a" help:"Show path to the script cache directory."`
	ShowState   bool `short:"s" help:"Show path to the script state directory."`
	ShowRuntime bool `short:"r" help:"Show path to the script runtime directory."`
}

func (c *CliCmdShow) Run(ctx cli.Context) error {
	switch {
	case c.Namespace.Value() == "":
		names, err := script.ListNamespaces("")
		if err != nil {
			return err
		}

		sort.Strings(names)

		for _, v := range names {
			fmt.Println(v)
		}

		return nil
	case c.Script == "":
		names, err := script.ListScripts(c.Namespace.Value(), "")
		if err != nil {
			return err
		}

		sort.Strings(names)

		for _, v := range names {
			fmt.Println(v)
		}

		return nil
	}

	scr, err := script.FindScript(c.Namespace.Value(), c.Script)
	if err != nil {
		return err
	}

	if err := scr.Init(); err != nil {
		return fmt.Errorf("cannot initialize script: %w", err)
	}

	defer scr.Cleanup()

	switch {
	case c.ShowPath, c.ShowConfig, c.ShowData, c.ShowCache, c.ShowState, c.ShowRuntime:
		c.showPaths(scr)

		return nil
	}

	return c.showPrettyData(ctx, scr)
}

func (c *CliCmdShow) showPrettyData(ctx context.Context, scr *script.Script) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

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

	c.showPrettyDataPaths(scr, dirSizes)
	fmt.Println()

	c.showPrettyDataGlobals(scr)
	fmt.Println()

	c.showParameters(scr)
	fmt.Println()

	c.showFlags(scr)

	return nil
}

func (c *CliCmdShow) showPaths(scr *script.Script) {
	if c.ShowPath {
		fmt.Println(scr.Path())
	}

	if c.ShowConfig {
		fmt.Println(scr.ConfigPath())
	}

	if c.ShowData {
		fmt.Println(scr.DataPath())
	}

	if c.ShowCache {
		fmt.Println(scr.CachePath())
	}

	if c.ShowState {
		fmt.Println(scr.StatePath())
	}

	if c.ShowRuntime {
		fmt.Println(scr.RuntimePath())
	}
}

func (c *CliCmdShow) showPrettyDataPaths(scr *script.Script, dirSizes map[string]*atomic.Int64) {
	table := tabby.New()

	table.AddLine("Path:", scr.Path(), "")
	table.AddLine("Config path:", scr.ConfigPath(), "")
	table.AddLine("Data path:", scr.DataPath(), c.getDirSize(dirSizes[scr.DataPath()].Load()))
	table.AddLine("Cache path:", scr.CachePath(), c.getDirSize(dirSizes[scr.CachePath()].Load()))
	table.AddLine("State path:", scr.StatePath(), c.getDirSize(dirSizes[scr.StatePath()].Load()))
	table.AddLine("Runtime path:", scr.RuntimePath(), c.getDirSize(dirSizes[scr.RuntimePath()].Load()))
	table.Print()
}

func (c *CliCmdShow) showPrettyDataGlobals(scr *script.Script) {
	table := tabby.New()
	conf := scr.Config()

	table.AddLine("Network:", strconv.FormatBool(conf.Network))
	table.AddLine("Git:", conf.Git.String())

	table.Print()
}

func (c *CliCmdShow) showParameters(scr *script.Script) {
	table := tabby.New()
	conf := scr.Config()

	if len(conf.Parameters) == 0 {
		return
	}

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
		valueI := c.boolToInt(paramI.Required())
		valueJ := c.boolToInt(paramJ.Required())

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
			c.showRequired(param.Required()),
			param.Type(),
			c.showParametersSpecification(param.Specification()))
	}

	table.Print()
}

func (c *CliCmdShow) showFlags(scr *script.Script) {
	table := tabby.New()
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
		valueI := c.boolToInt(flagI.Required())
		valueJ := c.boolToInt(flagJ.Required())

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
			c.showRequired(flag.Required()))
	}

	table.Print()
}

func (c *CliCmdShow) showRequired(isRequired bool) string {
	if isRequired {
		return RequiredTrue
	}

	return RequiredFalse
}

func (c *CliCmdShow) showParametersSpecification(spec map[string]string) string {
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

func (c *CliCmdShow) getDirSize(sizeInBytes int64) string {
	if sizeInBytes < 0 {
		return DirSizeError
	}

	size := float64(sizeInBytes)
	unit := 0

	for size >= ByteBase && unit < len(ByteUnits) {
		size /= ByteBase
		unit++
	}

	return strconv.FormatFloat(size, 'f', 2, 64) + ByteUnits[unit]
}

func (c *CliCmdShow) boolToInt(value bool) int {
	if value {
		return 1
	}

	return 0
}
