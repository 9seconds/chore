package edit

import (
	"embed"
	"text/template"
)

//go:embed static/*
var staticFS embed.FS

func getTemplate(path string) *template.Template {
	return template.Must(template.ParseFS(staticFS, path))
}
