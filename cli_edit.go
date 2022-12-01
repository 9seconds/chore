package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/9seconds/chore/internal/cli"
	"github.com/9seconds/chore/internal/script"
)

type CliCmdEdit struct {
	Editor cli.Editor `short:"e" help:"Editor to use."`

	Namespace cli.Namespace `arg:"" help:"Script namespace."`
	Script    string        `arg:"" help:"Script name."`
}

func (c *CliCmdEdit) Run(ctx cli.Context) error {
	executable, err := script.New(c.Namespace.Value(), c.Script)
	if err != nil {
		return fmt.Errorf("cannot initialize script: %w", err)
	}

	defer executable.Cleanup() //nolint: errcheck

	editor, err := c.Editor.Value()
	if err != nil {
		return fmt.Errorf("cannot initialize editor: %w", err)
	}

	cmd := exec.CommandContext(ctx, editor, executable.Path())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("editor didn't succeed: %w", err)
	}

	return nil
}
