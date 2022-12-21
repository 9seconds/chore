package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"

	"github.com/9seconds/chore/internal/cli"
	"github.com/alecthomas/kong"
)

var version = "dev"

var CLI struct {
	Debug   bool             `short:"d" env:"CHORE_DEBUG" help:"Run in debug mode."`
	Version kong.VersionFlag `short:"V" help:"Show version."`

	Show       CliCmdShow       `cmd:"" aliases:"s" help:"Show details on namespaces or scripts."`
	EditScript CliCmdEditScript `cmd:"" aliases:"e" help:"Edit chore script."`
	EditConfig CliCmdEditConfig `cmd:"" aliases:"c" help:"Edit chore script config."`
	Run        CliCmdRun        `cmd:"" aliases:"r" help:"Run chore script."`
}

type editorCommand struct {
	Editor cli.Editor `short:"e" help:"Editor to use."`

	Namespace cli.Namespace `arg:"" help:"Script namespace."`
	Script    string        `arg:"" help:"Script name."`
}

func (e *editorCommand) Open(ctx context.Context, path string, defaultContent []byte) error { //nolint: cyclop
	editor, err := e.Editor.Value()
	if err != nil {
		return fmt.Errorf("cannot initialize editor: %w", err)
	}

	var (
		templateContent = defaultContent
		originalContent []byte
		mode            fs.FileMode = 0o600
	)

	stat, err := os.Stat(path)

	switch {
	case err == nil:
		originalContent, err = os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("cannot read from original file: %w", err)
		}

		templateContent = originalContent
		mode = stat.Mode().Perm()
	case !errors.Is(err, fs.ErrNotExist):
		return fmt.Errorf("cannot stat original file: %w", err)
	}

	tempFile, err := os.CreateTemp("", "")
	if err != nil {
		return fmt.Errorf("cannot create temporary file: %w", err)
	}

	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write(templateContent); err != nil {
		return fmt.Errorf("cannot populate temporary file: %w", err)
	}

	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("cannot close temporary file: %w", err)
	}

	cmd := exec.CommandContext(ctx, editor, tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cannot successfully complete text editor: %w", err)
	}

	newContent, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return fmt.Errorf("cannot read an updated content of the file: %w", err)
	}

	if bytes.Equal(newContent, originalContent) {
		return nil
	}

	if err := os.WriteFile(path, newContent, mode); err != nil {
		return fmt.Errorf("cannot write content back to the original file: %w", err)
	}

	return nil
}
