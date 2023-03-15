package edit

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"

	"github.com/spf13/cobra"
)

type mainCallback func([]string, io.Writer) (string, fs.FileMode, error)

func main(callback mainCallback) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		editor, err := cmd.Flags().Lookup("editor").Value.(*FlagEditor).Get()
		if err != nil {
			return fmt.Errorf("cannot get editor: %w", err)
		}

		buf := bytes.Buffer{}

		path, mode, err := callback(args, &buf)
		if err != nil {
			return err
		}

		if err := ensureFile(path, buf.Bytes()); err != nil {
			return fmt.Errorf("cannot ensure file %s: %w", path, err)
		}

		if err := openEditor(cmd.Context(), editor, path); err != nil {
			return fmt.Errorf("cannot correctly finish editor: %w", err)
		}

		if err := os.Chmod(path, mode); err != nil {
			return fmt.Errorf("cannot set correct permissions to %s: %w", path, err)
		}

		return nil
	}
}
