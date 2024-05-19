package web

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed static/*
var staticFiles embed.FS

//go:embed html/*.html
var templateFiles embed.FS

func ServeStaticFiles() http.Handler {
	handler := http.FileServerFS(staticFiles)
	return handler
}

func ParseTemplates() *template.Template {
	templates, err := template.ParseFS(templateFiles, "html/*.html")
	if err != nil {
		panic(err)
	}

	return templates
}
