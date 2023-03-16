package completions

import (
	"log"
	"sort"

	"github.com/9seconds/chore/internal/script"
	"github.com/spf13/cobra"
)

func CompleteNamespaces(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	namespaces, err := script.ListNamespaces()
	if err != nil {
		log.Printf("cannot list namespaces: %v", err)

		return nil, cobra.ShellCompDirectiveError
	}

	return namespaces, cobra.ShellCompDirectiveNoFileComp
}

func CompleteAllNamespaceScripts(
	cmd *cobra.Command,
	args []string,
	toComplete string,
) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		return CompleteNamespaces(cmd, args, toComplete)
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

	toShow := map[string]bool{}

	for _, v := range scripts {
		toShow[v] = true
	}

	for _, v := range args[1:] {
		delete(toShow, v)
	}

	scripts = scripts[:0]

	for k := range toShow {
		scripts = append(scripts, k)
	}

	sort.Strings(scripts)

	return scripts, cobra.ShellCompDirectiveNoFileComp
}

func CompleteNamespaceScript(
	cmd *cobra.Command,
	args []string,
	toComplete string,
) ([]string, cobra.ShellCompDirective) {
	if len(args) < 2 { //nolint: gomnd
		return CompleteAllNamespaceScripts(cmd, args, toComplete)
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}
