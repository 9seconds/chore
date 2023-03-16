package vault

import (
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/9seconds/chore/internal/config"
	"github.com/9seconds/chore/internal/paths"
	"github.com/9seconds/chore/internal/script"
	"github.com/9seconds/chore/internal/vault"
	"github.com/spf13/cobra"
)

type mainCallback func(*cobra.Command, vault.Vault, []string) (bool, error)

func main(callback mainCallback) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		namespace, _ := script.ExtractRealNamespace(args[0])

		conf, err := config.Get()
		if err != nil {
			return fmt.Errorf("cannot get application config: %w", err)
		}

		if _, ok := conf.Vault[namespace]; !ok {
			return fmt.Errorf("cannot find out correct password for namespace %s", namespace)
		}

		vaultPath := paths.ConfigNamespaceScriptVault(namespace)

		vlt, err := getVault(vaultPath, conf.Vault[namespace])
		if err != nil {
			return fmt.Errorf("cannot open vault: %w", err)
		}

		save, err := callback(cmd, vlt, args[1:])

		switch {
		case err != nil:
			return err
		case save:
			return saveVault(vaultPath, vlt)
		}

		return nil
	}
}

func getVault(path, password string) (vault.Vault, error) {
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

func saveVault(path string, vlt vault.Vault) error {
	writer, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}

	defer writer.Close()

	return vault.Save(writer, vlt)
}
