package vault

import (
	"github.com/9seconds/chore/internal/cli/completions"
	"github.com/9seconds/chore/internal/cli/validators"
	"github.com/9seconds/chore/internal/vault"
	"github.com/spf13/cobra"
)

func NewSet() *cobra.Command {
	return &cobra.Command{
		Use:                   "set namespace key value",
		Short:                 "Set a vault secret",
		ValidArgsFunction:     completions.CompleteNamespaces,
		DisableFlagsInUseLine: true,
		Args: cobra.MatchAll(
			cobra.ExactArgs(3), //nolint: gomnd
			validators.Namespace(0),
		),
		RunE: main(func(cmd *cobra.Command, vlt vault.Vault, args []string) (bool, error) {
			vlt.Set(args[0], args[1])

			return true, nil
		}),
	}
}
