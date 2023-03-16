package cli

import (
	"github.com/9seconds/chore/internal/cli/edit"
	"github.com/spf13/cobra"
)

func NewEdit() *cobra.Command {
	cmd := &cobra.Command{
		Use:              "edit",
		Short:            "Edit scripts and configs",
		TraverseChildren: true,
	}

	var editorFlag edit.FlagEditor

	cmd.PersistentFlags().VarP(&editorFlag, "editor", "e", "editor to use")

	cmd.AddCommand(
		edit.NewScript(),
		edit.NewScriptConfig(),
		edit.NewAppConfig(),
	)

	return cmd
}
