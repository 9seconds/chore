package cli

import (
	"bytes"
	"fmt"
	"os"

	"github.com/9seconds/chore/internal/access"
	"github.com/9seconds/chore/internal/env"
	"github.com/9seconds/chore/internal/paths"
	"github.com/9seconds/chore/internal/script"
	"github.com/spf13/cobra"
)

func NewEditScript() *cobra.Command {
	cmd := &cobra.Command{
		Use:        "edit-script [flags] namespace script",
		Aliases:    []string{"e", "es"},
		SuggestFor: []string{"edit-config", "ec"},
		Short:      "Edit chore script",
		Args: cobra.MatchAll(
			cobra.ExactArgs(2), //nolint: gomnd
			validASCIIName(0, ErrNamespaceInvalid),
			validASCIIName(1, ErrScriptInvalid),
		),
		RunE:              mainEditScript,
		ValidArgsFunction: completeNamespaceScript,
	}

	var editorFlag flagEditor

	flags := cmd.Flags()

	flags.VarP(&editorFlag, "editor", "e", "editor to use")

	return cmd
}

func mainEditScript(cmd *cobra.Command, args []string) error {
	editor, err := cmd.Flag("editor").Value.(*flagEditor).Get()
	if err != nil {
		return err
	}

	scr := &script.Script{
		Namespace:  args[0],
		Executable: args[1],
	}

	defer scr.Cleanup()

	if err := script.EnsureDir(paths.ConfigNamespace(scr.Namespace)); err != nil {
		return fmt.Errorf("cannot ensure namespace dir: %w", err)
	}

	path := scr.Path()
	tpl := getTemplate("static/edit-script.sh")
	defaultContent := bytes.Buffer{}

	if err := tpl.Execute(&defaultContent, map[string]string{
		"DebugVar":     env.Debug,
		"DebugEnabled": env.DebugEnabled,
	}); err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	if err := openEditor(cmd.Context(), editor, path, defaultContent.Bytes()); err != nil {
		return fmt.Errorf("editor failed: %w", err)
	}

	if err := access.Access(path, false, false, true); err == nil {
		return nil
	}

	stat, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("cannot stat file: %w", err)
	}

	if err := os.Chmod(path, stat.Mode().Perm()|0o100); err != nil { //nolint: gomnd
		return fmt.Errorf("cannot set permissions: %w", err)
	}

	return nil
}
