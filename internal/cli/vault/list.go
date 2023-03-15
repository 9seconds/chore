package vault

import (
	"sort"

	"github.com/9seconds/chore/internal/cli/completions"
	"github.com/9seconds/chore/internal/cli/validators"
	"github.com/9seconds/chore/internal/vault"
	"github.com/spf13/cobra"
)

func NewList() *cobra.Command {
	return &cobra.Command{
		Use:                   "list namespace",
		Short:                 "List keys of a vault secrets",
		ValidArgsFunction:     completions.CompleteNamespaces,
		DisableFlagsInUseLine: true,
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), validators.Namespace(0)),
		RunE: main(func(cmd *cobra.Command, vlt vault.Vault, _ []string) (bool, error) {
			keys := vlt.List()

			sort.Strings(keys)

			for _, v := range keys {
				cmd.Println(v)
			}

			return false, nil
		}),
	}
}
