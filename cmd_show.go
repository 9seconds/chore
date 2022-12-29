package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/9seconds/chore/internal/cli"
	"github.com/9seconds/chore/internal/script"
)

type CliCmdShow struct {
	Namespace cli.Namespace `arg:"" optional:"" help:"Script namespace. Dot takes one from environment variable CHORE_NAMESPACE."`
	Script    string        `arg:"" optional:"" help:"Script name."`
}

func (c *CliCmdShow) Run(_ cli.Context) error {
	switch {
	case c.Namespace.Value() == "":
		names, err := script.ListNamespaces("")
		if err != nil {
			return err
		}

		sort.Strings(names)

		for _, v := range names {
			fmt.Println(v)
		}

		return nil
	case c.Script == "":
		names, err := script.ListScripts(c.Namespace.Value(), "")
		if err != nil {
			return err
		}

		sort.Strings(names)

		for _, v := range names {
			fmt.Println(v)
		}

		return nil
	}

	scr, err := script.FindScript(c.Namespace.Value(), c.Script)
	if err != nil {
		return err
	}

	if err := scr.Init(); err != nil {
		return fmt.Errorf("cannot initialize script: %w", err)
	}

	defer scr.Cleanup()

	return getTemplate("static/show.txt").Execute(os.Stdout, scr)
}
