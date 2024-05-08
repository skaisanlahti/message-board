package user

import (
	"bytes"
	"database/sql"
	"embed"
	"html/template"
	"net/http"
	"time"

	"github.com/skaisanlahti/message-board/internal/assert"
)

//go:embed html/*.html
var htmlFiles embed.FS

func RegisterPage(database *sql.DB) http.Handler {
	templates, err := template.ParseFS(htmlFiles, "html/*.html")
	assert.Ok(err, "failed to parse register templates")

	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		data := pageData{
			Key: time.Now().UnixMilli(),
		}

		html := renderPage(templates, data)
		response.WriteHeader(http.StatusOK)
		response.Write(html)
	})
}

type pageData struct {
	Key      int64
	Username string
	Password string
	Error    string
}

func renderPage(templates *template.Template, data pageData) []byte {
	var buffer bytes.Buffer
	err := templates.ExecuteTemplate(&buffer, "page", data)
	assert.Ok(err, "failed to render homepage")
	return buffer.Bytes()
}
