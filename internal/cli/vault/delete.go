package vault

import (
	"github.com/9seconds/chore/internal/cli/completions"
	"github.com/9seconds/chore/internal/cli/validators"
	"github.com/9seconds/chore/internal/vault"
	"github.com/spf13/cobra"
)

func NewDelete() *cobra.Command {
	return &cobra.Command{
		Use:                   "delete namespace key...",
		Short:                 "Delete vault secrets",
		ValidArgsFunction:     completions.CompleteNamespaces,
		DisableFlagsInUseLine: true,
		Args: cobra.MatchAll(
			cobra.MinimumNArgs(2), //nolint: gomnd
			validators.Namespace(0),
		),
		RunE: main(func(cmd *cobra.Command, vlt vault.Vault, args []string) (bool, error) {
			for _, v := range args {
				vlt.Delete(v)
			}

			return true, nil
		}),
	}
}
