package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/9seconds/chore/internal/argparse"
	"github.com/9seconds/chore/internal/cli"
	"github.com/9seconds/chore/internal/script"
)

type CliFishCompletion struct {
	Namespace cli.Namespace `arg:"" optional:"" help:"Script namespace."`
	Script    string        `arg:"" optional:"" help:"Script name."`
	Arguments []string      `arg:"" optional:"" passthrough:"" help:"Arguments that were passed"`
}

func (c *CliFishCompletion) Run(_ cli.Context) error { //nolint: unparam,cyclop
	if c.Namespace.Value() == "" {
		content, err := staticFS.ReadFile("static/fish-completion.fish")
		if err != nil {
			panic(err)
		}

		fmt.Println(string(content))

		return nil
	}

	scr, err := script.FindScript(c.Namespace.Value(), c.Script)
	if err != nil {
		log.Printf("cannot find out script %v", err)

		return nil
	}

	if err := scr.Init(); err != nil {
		log.Printf("cannot initialize script %v", err)

		return nil
	}

	defer scr.Cleanup()

	parsed, err := argparse.Parse(c.Arguments)

	switch {
	case err != nil:
		log.Printf("cannot parse arguments: %v", err)

		return nil
	case parsed.IsPositionalTime():
		log.Printf("positional time, skip completions")

		return nil
	}

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

	sort.Strings(completions)

	for _, v := range completions {
		fmt.Println(v)
	}

	return nil
}
