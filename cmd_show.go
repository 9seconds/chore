package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/9seconds/chore/internal/cli"
	"github.com/9seconds/chore/internal/script"
)

const (
	DirSizeUnknown = "?"
	DirSizeError   = "ERROR"

	ByteBase = 1024
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

type cliCmdShowPrettyContext struct {
	*script.Script

	dirSizes map[string]*atomic.Int64
}

func (c cliCmdShowPrettyContext) DirSize(path string) string {
	value := c.dirSizes[path]
	if value == nil {
		return DirSizeUnknown
	}

	sizeInBytes := value.Load()

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

	tplContext := cliCmdShowPrettyContext{
		Script: scr,
		dirSizes: map[string]*atomic.Int64{
			scr.DataPath():    {},
			scr.StatePath():   {},
			scr.CachePath():   {},
			scr.RuntimePath(): {},
		},
	}

	waiters := &sync.WaitGroup{}

	waiters.Add(len(tplContext.dirSizes))

	for path, accumulator := range tplContext.dirSizes {
		go c.calculateDirectorySize(ctx, waiters, path, accumulator)
	}

	waiters.Wait()

	return getTemplate("static/show.txt").Execute(os.Stdout, tplContext)
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

func (c *CliCmdShow) calculateDirectorySize(
	ctx context.Context,
	waiters *sync.WaitGroup,
	path string,
	accumulator *atomic.Int64,
) {
	defer waiters.Done()

	err := filepath.Walk(path, func(_ string, info fs.FileInfo, _ error) error {
		accumulator.Add(info.Size())

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			return nil
		}
	})
	if err != nil {
		log.Printf("cannot collect data size for %s: %v", path, err)
	}
}
