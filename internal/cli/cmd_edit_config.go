package cli

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/9seconds/chore/internal/config"
	"github.com/9seconds/chore/internal/paths"
	"github.com/9seconds/chore/internal/script"
	"github.com/spf13/cobra"
)

type editConfigTemplateContext struct {
	Vault map[string]string
}

func NewEditConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:        "edit-config",
		Aliases:    []string{"c"},
		SuggestFor: []string{},
		Short:      "Edit chore configuration TOML",
		Args:       cobra.MaximumNArgs(0),
		RunE:       mainEditConfig,
	}

	var editorFlag flagEditor

	cmd.Flags().VarP(&editorFlag, "editor", "e", "editor to use")

	return cmd
}

func mainEditConfig(cmd *cobra.Command, args []string) error {
	editor, err := cmd.Flag("editor").Value.(*flagEditor).Get()
	if err != nil {
		return err
	}

	namespaces, err := script.ListNamespaces()
	if err != nil {
		namespaces = nil
	}

	if len(namespaces) == 0 {
		namespaces = []string{"example_namespace"}
	}

	context := editConfigTemplateContext{
		Vault: make(map[string]string),
	}

	for _, ns := range namespaces {
		context.Vault[ns] = config.GeneratePassword()
	}

	path := paths.AppConfigPath()

	if err := script.EnsureDir(filepath.Dir(path)); err != nil {
		return fmt.Errorf("cannot ensure dir: %w", err)
	}

	defaultContent := bytes.Buffer{}
	tpl := getTemplate("static/edit-config-template.toml")

	if err := tpl.Execute(&defaultContent, context); err != nil {
		return fmt.Errorf("cannot render default template: %w", err)
	}

	return openEditor(cmd.Context(), editor, path, defaultContent.Bytes())
}
