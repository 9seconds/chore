package main

import (
	"fmt"
	"sort"

	"github.com/9seconds/chore/internal/cli"
	"github.com/9seconds/chore/internal/script"
)

type CliCmdList struct {
	Namespace cli.Namespace `arg:"" optional:"" help:"Namespace to list."`
}

func (c *CliCmdList) Run(_ cli.Context) error {
	var (
		names []string
		err   error
	)

	if c.Namespace.Value() == "" {
		names, err = script.ListNamespaces("")
	} else {
		names, err = script.ListScripts(c.Namespace.Value(), "")
	}

	if err != nil {
		return err
	}

	sort.Strings(names)

	for _, v := range names {
		fmt.Println(v)
	}

	return nil
}
