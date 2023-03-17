package vault

import (
	"errors"

	"github.com/9seconds/chore/internal/cli/validators"
	"github.com/9seconds/chore/internal/vault"
	"github.com/spf13/cobra"
)

var ErrKeyUnknown = errors.New("key is unknown")

func NewGet() *cobra.Command {
	return &cobra.Command{
		Use:                   "get namespace key",
		Aliases:               []string{"g"},
		Short:                 "Get vaule of a vault secret",
		ValidArgsFunction:     completeSecretKey,
		DisableFlagsInUseLine: true,
		Args: cobra.MatchAll(
			cobra.ExactArgs(2), //nolint: gomnd
			validators.Namespace(0),
		),
		RunE: main(func(cmd *cobra.Command, vlt vault.Vault, args []string) (bool, error) {
			value, ok := vlt.Get(args[0])

			if !ok {
				return false, ErrKeyUnknown
			}

			cmd.Println(value)

			return false, nil
		}),
	}
}
