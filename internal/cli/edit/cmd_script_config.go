package edit

import (
	"fmt"
	"io"
	"io/fs"

	"github.com/9seconds/chore/internal/cli/completions"
	"github.com/9seconds/chore/internal/cli/validators"
	"github.com/9seconds/chore/internal/paths"
	"github.com/9seconds/chore/internal/script"
	"github.com/spf13/cobra"
)

func NewScriptConfig() *cobra.Command {
	return &cobra.Command{
		Use:               "script-config namespace script",
		Aliases:           []string{"c"},
		Short:             "Edit chore script configuration TOML",
		ValidArgsFunction: completions.CompleteNamespaceScript,
		Args: cobra.MatchAll(
			cobra.ExactArgs(2), //nolint: gomnd
			validators.Script(0, 1),
		),
		Run: main(func(args []string, content io.Writer) (string, fs.FileMode, error) {
			namespace, _ := script.ExtractRealNamespace(args[0])
			scr := &script.Script{
				Namespace:  namespace,
				Executable: args[1],
			}

			if err := script.EnsureDir(paths.ConfigNamespace(namespace)); err != nil {
				return "", 0, fmt.Errorf("cannot ensure namespace dir: %w", err)
			}

			tpl := getTemplate("static/edit-script-config-template.toml")

			if err := tpl.Execute(content, scr); err != nil {
				return "", 0, fmt.Errorf("cannot render default template: %w", err)
			}

			return scr.ConfigPath(), ConfigDefaultPermission, nil
		}),
	}
}
