package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/9seconds/chore/internal/env"
	"github.com/9seconds/chore/internal/script"
	"github.com/adrg/xdg"
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
	choreDir := filepath.Join(xdg.ConfigHome, env.ChoreDir)
	entries, err := os.ReadDir(choreDir)
	if err != nil {
		return fmt.Errorf("cannot read chore dir %s: %w", choreDir, err)
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
	path := filepath.Join(xdg.ConfigHome, env.ChoreDir, c.Namespace.Value)

	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("cannot list scripts in namespace %s: %w", c.Namespace.Value, err)
	}

	names := make([]string, 0, len(entries))

	for _, v := range entries {
		if _, err := script.New(c.Namespace.Value, v.Name()); err != nil {
			names = append(names, v.Name())
		}
	}

	sort.Strings(names)

	for _, v := range names {
		fmt.Println(v)
	}

	return nil
}
