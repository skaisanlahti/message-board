package templ

import (
	"embed"
	"text/template"

	"github.com/skaisanlahti/message-board/internal/assert"
)

func ParseFS(fs embed.FS, pattern string) *template.Template {
	templates, err := template.ParseFS(fs, pattern)
	assert.Ok(err, "failed to parse html templates")
	assert.NotNil(templates, "templates were not initialized")
	return templates
}
