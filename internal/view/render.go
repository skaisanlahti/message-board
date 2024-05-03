package view

import (
	"bytes"
	"embed"
	"text/template"

	"github.com/skaisanlahti/message-board/internal/assert"
)

//go:embed html/*.html
var html embed.FS

var templates *template.Template

func init() {
	var err error
	templates, err = template.ParseFS(html, "html/*.html")
	assert.Ok(err, "failed to parse html templates")
	assert.NotNil(templates, "templates were not initialized")
}

func Render(name string, data any) []byte {
	html := &bytes.Buffer{}
	err := templates.ExecuteTemplate(html, name, data)
	assert.Ok(err, "failed to render template")
	assert.True(len(html.Bytes()) > 0, "rendered nothing")

	return html.Bytes()
}
