package view

import (
	"bytes"
	"embed"
	"text/template"
)

//go:embed html/*.html
var html embed.FS

var templates *template.Template

func init() {
	var err error
	templates, err = template.ParseFS(html, "html/*.html")
	if err != nil {
		panic(err)
	}
}

func Render(name string, data any) []byte {
	html := &bytes.Buffer{}
	err := templates.ExecuteTemplate(html, name, data)
	if err != nil {
		panic(err)
	}

	return html.Bytes()
}
