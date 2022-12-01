package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/9seconds/chore/internal/cli"
	"github.com/9seconds/chore/internal/script"
)

const (
	cliCmdEditConfigText = `{
  "description" : "Amazing {{ js .Executable }} of {{ js .Namespace }}!",
  "git": "if_undefined",
  "network": false,
  "parameters": {
    "param": {
      "type": "string",
      "required": false,
      "spec": {
      }
    }
  }
}`
)

var cliCmdEditConfigTemplate = template.Must(
	template.New("cli-cmd-edit-config").Parse(cliCmdEditConfigText))

type CliCmdEditConfig struct {
	editorCommand
}

func (c *CliCmdEditConfig) Run(ctx cli.Context) error {
	scr := &script.Script{
		Namespace:  c.Namespace.Value(),
		Executable: c.Script,
	}

	defer scr.Cleanup()

	if err := script.EnsureDir(scr.NamespacePath()); err != nil {
		return fmt.Errorf("cannot ensure namespace dir: %w", err)
	}

	path := scr.ConfigPath()

	if _, err := os.Stat(path); err != nil {
		file, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("cannot create new config %s: %w", path, err)
		}

		if err := cliCmdEditConfigTemplate.Execute(file, scr); err != nil {
			return err
		}

		file.Close()
	}

	if err := c.Open(ctx, path); err != nil {
		return fmt.Errorf("editor failed: %w", err)
	}

	return nil
}
