package cli

import (
	"github.com/9seconds/chore/internal/cli/vault"
	"github.com/spf13/cobra"
)

func NewVault() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "vault",
		Aliases: []string{"v"},
		Short:   "Access storage of secrets for the namespace",
	}

	rootCmd.AddCommand(
		vault.NewList(),
		vault.NewGet(),
		vault.NewSet(),
		vault.NewDelete())

	return rootCmd
}
