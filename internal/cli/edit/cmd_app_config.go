package edit

import (
	"fmt"
	"io"
	"io/fs"
	"sort"

	"github.com/9seconds/chore/internal/config"
	"github.com/9seconds/chore/internal/paths"
	"github.com/9seconds/chore/internal/script"
	"github.com/spf13/cobra"
)

type appConfigTemplateContext struct {
	Env   []string
	Vault map[string]string
}

func NewAppConfig() *cobra.Command {
	return &cobra.Command{
		Use:               "app-config",
		Aliases:           []string{"a"},
		Short:             "Edit chore configuration TOML",
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: main(func(args []string, content io.Writer) (string, fs.FileMode, error) {
			namespaces, err := script.ListNamespaces()
			if err != nil {
				namespaces = nil
			}

			if len(namespaces) == 0 {
				namespaces = []string{"example_namespace"}
			}

			context := appConfigTemplateContext{
				Vault: make(map[string]string),
			}

			for _, ns := range namespaces {
				context.Env = append(context.Env, ns)
				context.Vault[ns] = config.GeneratePassword()
			}

			sort.Strings(context.Env)

			tpl := getTemplate("static/edit-config-template.toml")
			if err := tpl.Execute(content, context); err != nil {
				return "", 0, fmt.Errorf("cannot render default template: %w", err)
			}

			return paths.AppConfigPath(), ConfigDefaultPermission, nil
		}),
	}
}
