package main

import (
	"fmt"

	"github.com/9seconds/chore/internal/cli"
	"github.com/9seconds/chore/internal/gc"
)

type CliCmdGC struct {
	DryRun bool `short:"n" help:"Dry run."`
}

func (c *CliCmdGC) Run(_ cli.Context) error {
	paths, err := gc.Collect()
	if err != nil {
		return fmt.Errorf("cannot collect paths: %w", err)
	}

	if !c.DryRun {
		return gc.Remove(paths)
	}

	for _, path := range paths {
		fmt.Println(path)
	}

	return nil
}
