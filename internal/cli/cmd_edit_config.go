package cli

import (
	"bytes"
	"fmt"

	"github.com/9seconds/chore/internal/paths"
	"github.com/9seconds/chore/internal/script"
	"github.com/spf13/cobra"
)

func NewEditConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:        "edit-config [flags] namespace script",
		Aliases:    []string{"c", "ec"},
		SuggestFor: []string{"edit-script", "es", "e"},
		Short:      "Edit chore script configuration TOML",
		Args: cobra.MatchAll(
			cobra.ExactArgs(2), //nolint: gomnd
			validASCIIName(0, ErrNamespaceInvalid),
			validASCIIName(1, ErrScriptInvalid),
		),
		RunE:              mainEditConfig,
		ValidArgsFunction: completeNamespaceScript,
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

	scr := &script.Script{
		Namespace:  args[0],
		Executable: args[1],
	}

	if err := script.EnsureDir(paths.ConfigNamespace(scr.Namespace)); err != nil {
		return fmt.Errorf("cannot ensure namespace dir: %w", err)
	}

	defaultContent := bytes.Buffer{}
	tpl := getTemplate("static/edit-config-template.toml")

	if err := tpl.Execute(&defaultContent, scr); err != nil {
		return fmt.Errorf("cannot render default template: %w", err)
	}

	return openEditor(
		cmd.Context(),
		editor,
		scr.ConfigPath(),
		defaultContent.Bytes())
}
