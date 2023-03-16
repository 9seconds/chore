package vault

import (
	"log"
	"sort"

	"github.com/9seconds/chore/internal/cli/completions"
	"github.com/9seconds/chore/internal/vault"
	"github.com/spf13/cobra"
)

func completeSecretKeys(
	cmd *cobra.Command,
	args []string,
	toComplete string,
) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return completions.CompleteNamespaces(cmd, args, toComplete)
	}

	keys := make(map[string]bool)
	toExecute := main(func(_ *cobra.Command, vlt vault.Vault, args []string) (bool, error) {
		for _, v := range vlt.List() {
			keys[v] = true
		}

		return false, nil
	})

	if err := toExecute(nil, args); err != nil {
		log.Printf("cannot get a list of keys: %v", err)

		return nil, cobra.ShellCompDirectiveError
	}

	for _, v := range args[1:] {
		delete(keys, v)
	}

	toShow := make([]string, 0, len(keys))

	for k := range keys {
		toShow = append(toShow, k)
	}

	sort.Strings(toShow)

	return toShow, cobra.ShellCompDirectiveNoFileComp
}

func completeSecretKey(
	cmd *cobra.Command,
	args []string,
	toComplete string,
) ([]string, cobra.ShellCompDirective) {
	if len(args) < 2 { //nolint: gomnd
		return completeSecretKeys(cmd, args, toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}
