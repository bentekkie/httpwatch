package web

import (
	"embed"
	"text/template"
)

var (

	//go:embed templates/*
	tmplFiles embed.FS

	// Template holds all the parsed templates in the templates directory
	Template *template.Template
)

func init() {
	var err error

	Template, err = template.ParseFS(tmplFiles, "templates/*.html")
	if err != nil {
		panic(err)
	}
}
