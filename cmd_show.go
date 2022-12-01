package main

import (
	"fmt"
	"os"
	"text/template"

	"github.com/9seconds/chore/internal/cli"
	"github.com/9seconds/chore/internal/script"
)

const (
	cliCmdShowText = `Path:           {{ .Path }}
Config path:    {{ .ConfigPath }}
Data path:      {{ .DataPath  }}
Cache path:     {{ .CachePath  }}
State path:     {{ .StatePath  }}
Runtime path:   {{ .RuntimePath  }}
Network:        {{ print .Config.Network }}
Git:            {{ print .Config.Git }}
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

var cliCmdShotTemplate = template.Must(
	template.New("show").Parse(cliCmdShowText))

type CliCmdShow struct {
	Namespace cli.Namespace `arg:"" help:"Script namespace. Dot takes one from environment variable CHORE_NAMESPACE."`
	Script    string        `arg:"" help:"Script name."`
}

func (c *CliCmdShow) Run(_ cli.Context) error {
	scr := &script.Script{
		Namespace:  c.Namespace.Value(),
		Executable: c.Script,
	}

	if err := scr.Init(); err != nil {
		return fmt.Errorf("cannot initialize script: %w", err)
	}

	defer scr.Cleanup()

	return cliCmdShotTemplate.Execute(os.Stdout, scr)
}
