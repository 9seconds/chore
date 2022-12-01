package main

import (
	"fmt"
	"os"

	"github.com/9seconds/chore/internal/access"
	"github.com/9seconds/chore/internal/cli"
	"github.com/9seconds/chore/internal/script"
)

type CliCmdEdit struct {
	editorCommand
}

func (c *CliCmdEdit) Run(ctx cli.Context) error {
	scr := script.Script{
		Namespace:  c.Namespace.Value(),
		Executable: c.Script,
	}

	defer scr.Cleanup()

	if err := script.EnsureDir(scr.NamespacePath()); err != nil {
		return fmt.Errorf("cannot ensure namespace dir: %w", err)
	}

	path := scr.Path()

	if err := c.Open(ctx, path); err != nil {
		return fmt.Errorf("editor failed: %w", err)
	}

	if err := access.Access(path, false, false, true); err == nil {
		return nil
	}

	stat, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("cannot stat script: %w", err)
	}

	if err := os.Chmod(path, stat.Mode().Perm()|0o100); err != nil {
		return fmt.Errorf("cannot set permissions: %w", err)
	}

	return nil
}
