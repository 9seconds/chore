package cli

import (
	"log"
	"sort"
	"strings"

	"github.com/9seconds/chore/internal/argparse"
	"github.com/9seconds/chore/internal/cli/completions"
	"github.com/9seconds/chore/internal/script"
	"github.com/spf13/cobra"
)

func completeRun(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) { //nolint: cyclop
	switch len(args) {
	case 0:
		return completions.CompleteNamespaces(cmd, args, toComplete)
	case 1:
		return completions.CompleteNamespaceScript(cmd, args, toComplete)
	}

	parsed, err := argparse.Parse(args[2:])

	switch {
	case err != nil:
		log.Printf("cannot parse arguments: %v", err)

		return nil, cobra.ShellCompDirectiveError
	case parsed.IsPositionalTime():
		return nil, cobra.ShellCompDirectiveDefault
	}

	scr, err := script.New(args[0], args[1])
	if err != nil {
		log.Printf("cannot initalize script: %v", err)

		return nil, cobra.ShellCompDirectiveError
	}

	completions := []string{}
	directive := cobra.ShellCompDirectiveNoFileComp

	for name, param := range scr.Config.Parameters {
		if _, ok := parsed.Parameters[name]; ok {
			continue
		}

		completion := name + string(argparse.SeparatorKeyword)

		if toComplete != "" {
			if !strings.HasPrefix(completion, toComplete) {
				continue
			}

			directive = cobra.ShellCompDirectiveNoSpace
		}

		if descr := param.Description(); descr != "" {
			completion += "\t" + descr
		}

		completions = append(completions, completion)
	}

	for name, flag := range scr.Config.Flags {
		if _, ok := parsed.Flags[name]; ok {
			continue
		}

		negative := string(argparse.PrefixFlagNegative) + name
		positive := string(argparse.PrefixFlagPositive) + name

		if descr := flag.Description(); descr != "" {
			negative += "\t" + descr + " (no)"
			positive += "\t" + descr + " (yes)"
		}

		if toComplete == "" || strings.HasPrefix(negative, toComplete) {
			completions = append(completions, negative)
		}

		if toComplete == "" || strings.HasPrefix(positive, toComplete) {
			completions = append(completions, positive)
		}
	}

	sort.Strings(completions)

	if len(completions) > 0 {
		return completions, directive
	}

	return nil, cobra.ShellCompDirectiveDefault
}
