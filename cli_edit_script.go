package main

import (
	"fmt"
	"os"

	"github.com/9seconds/chore/internal/access"
	"github.com/9seconds/chore/internal/cli"
	"github.com/9seconds/chore/internal/script"
)

type CliCmdEditScript struct {
	editorCommand
}

func (c *CliCmdEditScript) Run(ctx cli.Context) error {
	scr := &script.Script{
		Namespace:  c.Namespace.Value(),
		Executable: c.Script,
	}

	defer scr.Cleanup()

	if err := script.EnsureDir(scr.NamespacePath()); err != nil {
		return fmt.Errorf("cannot ensure namespace dir: %w", err)
	}

	path := scr.Path()

	defaultContent, err := staticFS.ReadFile("static/edit-script.sh")
	if err != nil {
		panic(err)
	}

	if err := c.Open(ctx, path, defaultContent); err != nil {
		return fmt.Errorf("editor failed: %w", err)
	}

	if err := access.Access(path, false, false, true); err == nil {
		return nil
	}

	stat, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("cannot stat file: %w", err)
	}

	if err := os.Chmod(path, stat.Mode().Perm()|0o100); err != nil { //nolint: gomnd
		return fmt.Errorf("cannot set permissions: %w", err)
	}

	return nil
}
