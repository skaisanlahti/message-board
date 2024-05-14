package web

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/skaisanlahti/message-board/internal/assert"
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
	assert.Ok(err, "failed to parse html templates")
	return templates
}
