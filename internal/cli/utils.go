package cli

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
)

const (
	fileDefaultPermission = 0o600
)

func openEditor(ctx context.Context, editor, path string, templateContent []byte) error {
	_, err := os.Stat(path)

	switch {
	case errors.Is(err, fs.ErrNotExist):
		if err := os.WriteFile(path, templateContent, fileDefaultPermission); err != nil {
			return fmt.Errorf("cannot populate file with a template content: %w", err)
		}
	case err != nil:
		return fmt.Errorf("cannot stat file: %w", err)
	}

	cmd := exec.CommandContext(ctx, editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
