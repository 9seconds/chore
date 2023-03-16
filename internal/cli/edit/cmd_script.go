package edit

import (
	"fmt"
	"io"
	"io/fs"

	"github.com/9seconds/chore/internal/cli/completions"
	"github.com/9seconds/chore/internal/cli/validators"
	"github.com/9seconds/chore/internal/env"
	"github.com/9seconds/chore/internal/paths"
	"github.com/9seconds/chore/internal/script"
	"github.com/spf13/cobra"
)

type scriptTemplateContext struct {
	DebugVar     string
	DebugEnabled string
}

func NewScript() *cobra.Command {
	return &cobra.Command{
		Use:               "script namespace script",
		Short:             "Edit script",
		ValidArgsFunction: completions.CompleteNamespaceScript,
		Args: cobra.MatchAll(
			cobra.ExactArgs(2), //nolint: gomnd
			validators.ASCIIName(0, validators.ErrNamespaceInvalid),
			validators.ASCIIName(1, validators.ErrScriptInvalid),
		),
		RunE: main(func(args []string, content io.Writer) (string, fs.FileMode, error) {
			namespace, _ := script.ExtractRealNamespace(args[0])
			scr := &script.Script{
				Namespace:  namespace,
				Executable: args[1],
			}

			if err := script.EnsureDir(paths.ConfigNamespace(namespace)); err != nil {
				return "", 0, fmt.Errorf("cannot ensure namespace dir: %w", err)
			}

			tpl := getTemplate("static/edit-script.sh")
			context := scriptTemplateContext{
				DebugVar:     env.Debug,
				DebugEnabled: env.DebugEnabled,
			}

			if err := tpl.Execute(content, context); err != nil {
				return "", 0, fmt.Errorf("cannot render default template: %w", err)
			}

			return scr.Path(), ScriptDefaultPermissions, nil
		}),
	}
}
