package home

import (
	"bytes"
	"database/sql"
	"embed"
	"html/template"
	"net/http"

	"github.com/skaisanlahti/message-board/internal/assert"
)

//go:embed html/*.html
var htmlFiles embed.FS

func HomePage(database *sql.DB) http.Handler {
	templates, err := template.ParseFS(htmlFiles, "html/*.html")
	assert.Ok(err, "failed to parse home templates")

	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		data := homePageData{
			Greeting: "hello world",
		}

		html := renderHomePage(templates, data)
		response.WriteHeader(http.StatusOK)
		response.Write(html)
	})
}

type homePageData struct {
	Greeting string
}

func renderHomePage(templates *template.Template, data homePageData) []byte {
	var buffer bytes.Buffer
	err := templates.ExecuteTemplate(&buffer, "home", data)
	assert.Ok(err, "failed to render homepage")
	return buffer.Bytes()
}
