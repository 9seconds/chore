package main

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/9seconds/chore/internal/cli"
	"github.com/9seconds/chore/internal/script"
)

const (
	cliCmdEditConfigText = `// this is HJSON: https://hjson.github.io/

{
  description : Amazing {{ js .Executable }} of {{ js .Namespace }}!

  // valid options are 'if_undefined', 'always' and 'no'
  git: if_undefined
  network: false

  // to run command as a user
  // as_user: root

  parameters: {
    param: {
      type: string
      required: false
      spec: {
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

	defaultContent := bytes.Buffer{}

	if err := cliCmdEditConfigTemplate.Execute(&defaultContent, scr); err != nil {
		return fmt.Errorf("cannot render default template: %w", err)
	}

	if err := c.Open(ctx, scr.Path(), defaultContent.Bytes()); err != nil {
		return fmt.Errorf("editor failed: %w", err)
	}

	return nil
}
