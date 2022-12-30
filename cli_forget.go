package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/9seconds/chore/internal/cli"
	"github.com/9seconds/chore/internal/script"
)

type CliCmdForget struct {
	Namespace cli.Namespace `arg:"" help:"Script namespace."`
	Script    string        `arg:"" help:"Script name."`

	CleanConfig bool `short:"g" help:"Remove config file."`
	KeepData    bool `short:"t" help:"Keep data."`
	KeepCache   bool `short:"c" help:"Keep cache."`
	KeepState   bool `short:"s" help:"Keep state."`
	KeepRuntime bool `short:"r" help:"Keep runtime."`
}

func (c *CliCmdForget) Run(_ cli.Context) error {
	scr := &script.Script{
		Namespace:  c.Namespace.Value(),
		Executable: c.Script,
	}

	defer scr.Cleanup()

	if err := script.ValidateScript(scr.Path()); err != nil {
		return fmt.Errorf("invalid script: %w", err)
	}

	if c.CleanConfig {
		if err := os.Remove(scr.ConfigPath()); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("cannot remove config: %w", err)
		}
	}

	if !c.KeepData {
		if err := os.RemoveAll(scr.DataPath()); err != nil {
			return fmt.Errorf("cannot remove data path: %w", err)
		}
	}

	if !c.KeepCache {
		if err := os.RemoveAll(scr.CachePath()); err != nil {
			return fmt.Errorf("cannot remove cache path: %w", err)
		}
	}

	if !c.KeepState {
		if err := os.RemoveAll(scr.StatePath()); err != nil {
			return fmt.Errorf("cannot remove state path: %w", err)
		}
	}

	if !c.KeepRuntime {
		if err := os.RemoveAll(scr.RuntimePath()); err != nil {
			return fmt.Errorf("cannot remove runtime path: %w", err)
		}
	}

	return nil
}
