package cli

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"sort"

	"github.com/9seconds/chore/internal/config"
	"github.com/9seconds/chore/internal/paths"
	"github.com/9seconds/chore/internal/vault"
	"github.com/spf13/cobra"
)

func NewVault() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "vault",
		Short: "Access storage of secrets for the namespace",
	}

	listCmd := &cobra.Command{
		Use:                   "list namespace",
		Short:                 "List keys of a vault secrets",
		ValidArgsFunction:     completeNamespace,
		DisableFlagsInUseLine: true,
		Args: cobra.MatchAll(
			cobra.ExactArgs(1),
			validNamespace(0),
		),
		RunE: mainVault(func(cmd *cobra.Command, vlt vault.Vault, _ []string) bool {
			keys := vlt.List()

			sort.Strings(keys)

			for _, v := range keys {
				cmd.Println(v)
			}

			return false
		}),
	}

	getCmd := &cobra.Command{
		Use:                   "get namespace key",
		Short:                 "Get a value from a vault",
		ValidArgsFunction:     completeNamespace,
		DisableFlagsInUseLine: true,
		Args: cobra.MatchAll(
			cobra.ExactArgs(2),
			validNamespace(0),
		),
		RunE: mainVault(func(cmd *cobra.Command, vlt vault.Vault, args []string) bool {
			if value, ok := vlt.Get(args[0]); ok {
				cmd.Println(value)
			}

			return false
		}),
	}

	setCmd := &cobra.Command{
		Use:                   "set namespace key value",
		Short:                 "Set a value to a vault",
		ValidArgsFunction:     completeNamespace,
		DisableFlagsInUseLine: true,
		Args: cobra.MatchAll(
			cobra.ExactArgs(3),
			validNamespace(0),
		),
		RunE: mainVault(func(cmd *cobra.Command, vlt vault.Vault, args []string) bool {
			vlt.Set(args[0], args[1])

			return true
		}),
	}

	delCmd := &cobra.Command{
		Use:                   "delete namespace key",
		Short:                 "Delete a key from a vault",
		ValidArgsFunction:     completeNamespace,
		DisableFlagsInUseLine: true,
		Args: cobra.MatchAll(
			cobra.ExactArgs(2),
			validNamespace(0),
		),
		RunE: mainVault(func(cmd *cobra.Command, vlt vault.Vault, args []string) bool {
			vlt.Delete(args[0])

			return true
		}),
	}

	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(delCmd)

	return rootCmd
}

func mainVault(callback func(*cobra.Command, vault.Vault, []string) bool) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		namespace, exists := extractRealNamespace(args[0])
		if !exists {
			return ErrNamespaceInvalid
		}

		conf, err := mainVaultGetConfig()
		if err != nil {
			return fmt.Errorf("cannot get application config: %w", err)
		}

		if _, ok := conf.Vault[namespace]; !ok {
			return fmt.Errorf("cannot find out correct password for namespace %s", namespace)
		}

		vaultPath := paths.ConfigNamespaceScriptVault(namespace)

		vlt, err := mainVaultGetVault(vaultPath, conf.Vault[namespace])
		if err != nil {
			return fmt.Errorf("cannot open vault: %w", err)
		}

		if callback(cmd, vlt, args[1:]) {
			return mainVaultSaveVault(vaultPath, vlt)
		}

		return nil
	}
}

func mainVaultGetConfig() (config.Config, error) {
	conf := config.Config{}

	confReader, err := os.Open(paths.AppConfigPath())

	switch {
	case errors.Is(err, fs.ErrNotExist):
		return conf, nil
	case err != nil:
		return conf, fmt.Errorf("cannot open application config: %w", err)
	}

	defer confReader.Close()

	return config.ReadConfig(confReader)
}

func mainVaultGetVault(path, password string) (vault.Vault, error) {
	reader, err := os.Open(path)

	switch {
	case errors.Is(err, fs.ErrNotExist):
		return vault.New(password)
	case err != nil:
		return nil, fmt.Errorf("cannot open vault: %w", err)
	}

	defer reader.Close()

	return vault.Open(reader, password)
}

func mainVaultSaveVault(path string, vlt vault.Vault) error {
	writer, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}

	defer writer.Close()

	return vault.Save(writer, vlt)
}
