package web

import (
	"embed"
	"text/template"
)

var (

	//go:embed templates/*
	tmplFiles embed.FS
	Template  *template.Template
)

func init() {
	var err error

	Template, err = template.ParseFS(tmplFiles, "templates/*.html")
	if err != nil {
		panic(err)
	}
}
