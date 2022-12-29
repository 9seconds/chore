package main

import (
	"bytes"
	"fmt"

	"github.com/9seconds/chore/internal/cli"
	"github.com/9seconds/chore/internal/script"
)

type CliCmdEditConfig struct {
	editorCommand
}

func (c *CliCmdEditConfig) Run(ctx cli.Context) error {
	tpl := getTemplate("static/edit-config-template.hjson")

	scr := &script.Script{
		Namespace:  c.Namespace.Value(),
		Executable: c.Script,
	}

	defer scr.Cleanup()

	if err := script.EnsureDir(scr.NamespacePath()); err != nil {
		return fmt.Errorf("cannot ensure namespace dir: %w", err)
	}

	defaultContent := bytes.Buffer{}

	if err := tpl.Execute(&defaultContent, scr); err != nil {
		return fmt.Errorf("cannot render default template: %w", err)
	}

	if err := c.Open(ctx, scr.ConfigPath(), defaultContent.Bytes()); err != nil {
		return fmt.Errorf("editor failed: %w", err)
	}

	_, err := c.RemoveIfEmpty(scr.ConfigPath())

	return err
}
