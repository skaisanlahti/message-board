package view

import (
	"bytes"
	"embed"
	"text/template"

	"github.com/skaisanlahti/message-board/internal/assert"
)

//go:embed html/*.html
var html embed.FS

type Renderer struct {
	templates *template.Template
}

func (r *Renderer) Render(name string, data any) []byte {
	html := &bytes.Buffer{}
	err := r.templates.ExecuteTemplate(html, name, data)
	assert.Ok(err, "failed to render template")
	assert.True(len(html.Bytes()) > 0, "rendered nothing")
	return html.Bytes()
}

func New() *Renderer {
	templates, err := template.ParseFS(html, "html/*.html")
	assert.Ok(err, "failed to parse html templates")
	assert.NotNil(templates, "templates were not initialized")
	return &Renderer{templates}
}
