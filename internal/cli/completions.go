package cli

import (
	"log"

	"github.com/9seconds/chore/internal/argparse"
	"github.com/9seconds/chore/internal/script"
	"github.com/spf13/cobra"
)

func completeNamespaces() ([]string, cobra.ShellCompDirective) {
	namespaces, err := script.ListNamespaces()
	if err != nil {
		log.Printf("cannot list namespaces: %v", err)

		return nil, cobra.ShellCompDirectiveError
	}

	return namespaces, cobra.ShellCompDirectiveNoFileComp
}

func completeScripts(namespace string) ([]string, cobra.ShellCompDirective) {
	scripts, err := script.ListScripts(namespace)
	if err != nil {
		log.Printf("cannot list scripts: %v", err)

		return nil, cobra.ShellCompDirectiveError
	}

	return scripts, cobra.ShellCompDirectiveNoFileComp
}

func completeNamespaceScript(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	switch len(args) {
	case 0:
		return completeNamespaces()
	case 1:
		return completeScripts(args[0])
	}

	return nil, cobra.ShellCompDirectiveNoFileComp
}

func completeRun(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) { //nolint: cyclop
	switch len(args) {
	case 0:
		return completeNamespaces()
	case 1:
		return completeScripts(args[0])
	}

	parsed, err := argparse.Parse(args[2:])

	switch {
	case err != nil:
		log.Printf("cannot parse arguments: %v", err)

		return nil, cobra.ShellCompDirectiveError
	case parsed.IsPositionalTime():
		return nil, cobra.ShellCompDirectiveDefault
	}

	scr := &script.Script{
		Namespace:  args[0],
		Executable: args[1],
	}

	if err := scr.Init(); err != nil {
		log.Printf("cannot initalize script: %v", err)

		return nil, cobra.ShellCompDirectiveError
	}

	defer scr.Cleanup()

	conf := scr.Config()
	completions := []string{}

	for name, param := range conf.Parameters {
		if _, ok := parsed.Parameters[name]; ok {
			continue
		}

		completion := name + string(argparse.SeparatorKeyword)

		if descr := param.Description(); descr != "" {
			completion += "\t" + descr
		}

		completions = append(completions, completion)
	}

	for name, flag := range conf.Flags {
		if _, ok := parsed.Flags[name]; ok {
			continue
		}

		negative := string(argparse.PrefixFlagNegative) + name
		positive := string(argparse.PrefixFlagPositive) + name

		if descr := flag.Description(); descr != "" {
			negative += "\t" + descr + " (no)"
			positive += "\t" + descr + " (yes)"
		}

		completions = append(completions, negative, positive)
	}

	if len(completions) > 0 {
		return completions, cobra.ShellCompDirectiveNoFileComp
	}

	return nil, cobra.ShellCompDirectiveDefault
}
