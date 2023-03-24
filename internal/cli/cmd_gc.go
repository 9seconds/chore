package cli

import (
	"fmt"

	"github.com/9seconds/chore/internal/gc"
	"github.com/9seconds/chore/internal/script"
	"github.com/spf13/cobra"
)

func NewGC() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "gc",
		Short:             "Cleanup garbage in chore directories",
		Args:              cobra.NoArgs,
		RunE:              mainGC,
		ValidArgsFunction: cobra.NoFileCompletions,
	}

	cmd.Flags().BoolP("dry-run", "n", false, "dry run")

	return cmd
}

func mainGC(cmd *cobra.Command, _ []string) error {
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	validScripts := []*script.Script{}

	namespaces, err := script.ListNamespaces()
	if err != nil {
		return fmt.Errorf("cannot get list of all namespaces: %w", err)
	}

	for _, namespace := range namespaces {
		scripts, err := script.ListScripts(namespace)
		if err != nil {
			return fmt.Errorf("cannot list scripts in %s: %w", namespace, err)
		}

		for _, name := range scripts {
			validScripts = append(validScripts, &script.Script{
				Namespace:  namespace,
				Executable: name,
			})
		}
	}

	paths, err := gc.Collect(validScripts)
	if err != nil {
		return fmt.Errorf("cannot collect paths: %w", err)
	}

	if !dryRun {
		return gc.Remove(paths)
	}

	for _, path := range paths {
		cmd.Println(path)
	}

	return nil
}
