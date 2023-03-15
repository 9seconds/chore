package edit

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/9seconds/chore/internal/commands"
)

const (
	ConfigDefaultPermission  = 0o666
	ScriptDefaultPermissions = 0o755
	DirDefaultPermissions    = 0o777
)

func ensureFile(path string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), DirDefaultPermissions); err != nil {
		return fmt.Errorf("cannot ensure parent dir %s: %w", filepath.Dir(path), err)
	}

	writer, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, ConfigDefaultPermission)

	switch {
	case errors.Is(err, fs.ErrExist):
		return nil
	case err != nil:
		return fmt.Errorf("cannot create a file: %w", err)
	}

	_, err = writer.Write(content)

	return err
}

func openEditor(ctx context.Context, editor, path string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

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
