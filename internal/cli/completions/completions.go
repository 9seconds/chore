package completions

import (
	"log"

	"github.com/9seconds/chore/internal/script"
	"github.com/spf13/cobra"
)

func CompleteNamespaces(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	namespaces, err := script.ListNamespaces()
	if err != nil {
		log.Printf("cannot list namespaces: %v", err)

		return nil, cobra.ShellCompDirectiveError
	}

	return namespaces, cobra.ShellCompDirectiveNoFileComp
}

func CompleteNamespaceScript(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	switch len(args) {
	case 0:
		return CompleteNamespaces(cmd, args, toComplete)
	case 1:
	default:
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	namespace, exists := script.ExtractRealNamespace(args[0])
	if !exists {
		log.Printf("namespace is not defined")

		return nil, cobra.ShellCompDirectiveError
	}

	scripts, err := script.ListScripts(namespace)
	if err != nil {
		log.Printf("cannot list scripts: %v", err)

		return nil, cobra.ShellCompDirectiveError
	}

	return scripts, cobra.ShellCompDirectiveNoFileComp
}
