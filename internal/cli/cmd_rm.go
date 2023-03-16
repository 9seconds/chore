package cli

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"sort"

	"github.com/9seconds/chore/internal/cli/completions"
	"github.com/9seconds/chore/internal/cli/validators"
	"github.com/9seconds/chore/internal/script"
	"github.com/spf13/cobra"
)

func NewRm() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "rm namespace script...",
		Short:                 "Remove scripts and its directories from a namespace",
		DisableFlagsInUseLine: true,
		ValidArgsFunction:     completions.CompleteAllNamespaceScripts,
		Args: cobra.MatchAll(
			cobra.MinimumNArgs(2),
			validators.Namespace(0),
		),
		RunE: mainRm,
	}

	cmd.Flags().BoolP("dry-run", "n", false, "dry run")

	return cmd
}

func mainRm(cmd *cobra.Command, args []string) error {
	dryRun, err := cmd.Flags().GetBool("dry-run")
	if err != nil {
		return err
	}

	namespace, _ := script.ExtractRealNamespace(args[0])
	toRemove := make([]string, 0, 5*(len(args)-1))

	for _, name := range args[1:] {
		scr := &script.Script{
			Namespace:  namespace,
			Executable: name,
		}
		toRemove = append(
			toRemove,
			scr.Path(),
			scr.ConfigPath(),
			scr.DataPath(),
			scr.CachePath(),
			scr.StatePath())
	}

	sort.Strings(toRemove)

	for _, path := range toRemove {
		if !dryRun {
			if err := os.RemoveAll(path); err != nil {
				return fmt.Errorf("cannot remove path %s: %w", path, err)
			}

			continue
		}

		_, err := os.Stat(path)

		switch {
		case err == nil:
			cmd.Println(path)
		case !errors.Is(err, fs.ErrNotExist):
			return fmt.Errorf("cannot stat path %s: %w", path, err)
		}
	}

	return nil
}
