package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/9seconds/chore/chorelib/env"
	"github.com/9seconds/chore/chorelib/script"
)

type CliCmdList struct {
	Namespace CliNamespace `arg:"" optional:"" help:"Namespace to list."`
}

func (c *CliCmdList) Run(ctx Context) error {
	if c.Namespace.Value == "" {
		return c.listNamespaces()
	}

	return c.listScripts()
}

func (c *CliCmdList) listNamespaces() error {
	entries, err := os.ReadDir(env.Home)
	if err != nil {
		return fmt.Errorf("cannot read home %s: %w", env.Home, err)
	}

	names := make([]string, 0, len(entries))

	for _, v := range entries {
		if v.IsDir() {
			names = append(names, v.Name())
		}
	}

	sort.Strings(names)

	for _, v := range names {
		fmt.Println(v)
	}

	return nil
}

func (c *CliCmdList) listScripts() error {
	path := filepath.Join(env.Home, c.Namespace.Value)

	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("cannot list scripts in namespace %s: %w", c.Namespace.Value, err)
	}

	names := make([]string, 0, len(entries))

	for _, v := range entries {
		vv := script.Script{
			Namespace:  c.Namespace.Value,
			Executable: v.Name(),
		}

		if vv.IsValid() != nil {
			names = append(names, vv.String())
		}
	}

	sort.Strings(names)

	for _, v := range names {
		fmt.Println(v)
	}

	return nil
}
