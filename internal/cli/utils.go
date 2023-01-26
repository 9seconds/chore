package cli

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/9seconds/chore/internal/commands"
)

const (
	fileDefaultPermission = 0o600
)

var ErrExpectedFile = errors.New("file is expected")

func openEditor(ctx context.Context, editor, path string, templateContent []byte) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	stat, err := os.Stat(path)

	switch {
	case errors.Is(err, fs.ErrNotExist):
		if err := os.WriteFile(path, templateContent, fileDefaultPermission); err != nil {
			return fmt.Errorf("cannot populate file with a template content: %w", err)
		}
	case stat.IsDir():
		return ErrExpectedFile
	case err != nil:
		return fmt.Errorf("cannot stat file: %w", err)
	}

	cmd := commands.New(editor, []string{path}, nil, os.Stdin, os.Stdout, os.Stderr)

	if err := cmd.Start(ctx); err != nil {
		return fmt.Errorf("cannot start editor: %w", err)
	}

	result := cmd.Wait()

	if result.ExitCode != 0 {
		return fmt.Errorf("command exited with %d", result.ExitCode)
	}

	return nil
}
