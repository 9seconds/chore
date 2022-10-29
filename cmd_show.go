package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/9seconds/chore/chorelib/script"
)

const (
	cliCmdShowText = `Path:           {{ .Path }}
Persistent dir: {{ .PersistentDir  }}
Network:        {{ print .Config.Network }}
{{ if .Config.Description }}
{{ .Config.Description }}
{{ end }}

{{- if .Config.Parameters }}
Parameters:
{{ range $key, $value := .Config.Parameters -}}
	{{- $key }} ({{ $value.Type }}) -> {{ $value }}
{{ end -}}

{{- end -}}`
)

var cliCmdShotTemplate = template.Must(template.New("show").Parse(cliCmdShowText))

type CliCmdShow struct {
	Namespace CliNamespace `arg:"" help:"Script namespace. Dot takes one from environment variable CHORE_NAMESPACE."`
	Script    string       `arg:"" help:"Script name."`
}

func (c *CliCmdShow) Run(ctx Context) error {
	executable := script.Script{
		Namespace:  c.Namespace.Value,
		Executable: c.Script,
	}

	if err := executable.IsValid(); err != nil {
		return fmt.Errorf("script is invalid: %w", err)
	}

	if err := executable.Init(); err != nil {
		return fmt.Errorf("cannot initialize script: %w", err)
	}

	defer executable.Cleanup()

	cliCmdShotTemplate.Execute(os.Stdout, executable)

	return nil
}
