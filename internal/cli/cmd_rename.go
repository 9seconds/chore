package cli

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/9seconds/chore/internal/cli/completions"
	"github.com/9seconds/chore/internal/cli/validators"
	"github.com/9seconds/chore/internal/script"
	"github.com/spf13/cobra"
)

func NewRename() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "rename [flags] namespace from to",
		Aliases:               []string{"mv"},
		Short:                 "Renames scripts and its directories.",
		DisableFlagsInUseLine: true,
		ValidArgsFunction:     completions.CompleteNamespaceScript,
		Args: cobra.MatchAll(
			cobra.ExactArgs(3), //nolint: gomnd
			validators.Script(0, 1),
		),
		RunE: mainRename,
	}

	cmd.Flags().BoolP("force", "f", false, "Do, not ask")

	return cmd
}

func mainRename(cmd *cobra.Command, args []string) error {
	if args[1] == args[2] {
		return nil
	}

	force, err := cmd.Flags().GetBool("force")
	if err != nil {
		return err
	}

	namespace, _ := script.ExtractRealNamespace(args[0])
	scriptFrom := &script.Script{
		Namespace:  namespace,
		Executable: args[1],
	}
	scriptTo := &script.Script{
		Namespace:  namespace,
		Executable: args[2],
	}

	moveMap := map[string]string{
		scriptFrom.Path():       scriptTo.Path(),
		scriptFrom.ConfigPath(): scriptTo.ConfigPath(),
		scriptFrom.DataPath():   scriptTo.DataPath(),
		scriptFrom.CachePath():  scriptTo.CachePath(),
		scriptFrom.StatePath():  scriptTo.StatePath(),
	}

	if err := mainRenameValidate(force, moveMap); err != nil {
		return err
	}

	return MainRenameProcess(moveMap)
}

func mainRenameValidate(force bool, moveMap map[string]string) error {
	for _, path := range moveMap {
		_, err := os.Stat(path)

		switch {
		case errors.Is(err, fs.ErrNotExist):
		case err != nil:
			return fmt.Errorf("cannot stat path %s: %w", path, err)
		case !force:
			return fmt.Errorf("path %s exists, it prevents renaming", path)
		}
	}

	return nil
}

func MainRenameProcess(moveMap map[string]string) error {
	for src, dest := range moveMap {
		if err := os.RemoveAll(dest); err != nil {
			return fmt.Errorf("cannot remove %s: %w", dest, err)
		}

		_, err := os.Stat(src)

		switch {
		case errors.Is(err, fs.ErrNotExist):
			continue
		case err != nil:
			return fmt.Errorf("cannot stat path: %w", err)
		}

		if err := os.Rename(src, dest); err != nil {
			return fmt.Errorf("cannot rename %s to %s: %w", src, dest, err)
		}
	}

	return nil
}
