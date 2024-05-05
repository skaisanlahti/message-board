package home

import (
	"bytes"
	"context"
	"database/sql"
	"embed"
	"fmt"
	"net/http"
	"text/template"

	"github.com/skaisanlahti/message-board/internal/assert"
	"github.com/skaisanlahti/message-board/internal/templ"
)

//go:embed html/*.html
var html embed.FS

func HomePage(database *sql.DB) http.Handler {
	templates := templ.ParseFS(html, "html/*.html")

	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		username, err := getUserName(database, request.Context(), "1234")
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			return
		}

		greeting := fmt.Sprintf("Hello %s", username)
		data := HomePageData{
			Greeting: greeting,
		}

		html := renderHomePage(templates, data)
		response.WriteHeader(http.StatusOK)
		response.Write(html)
	})
}

func getUserName(database *sql.DB, ctx context.Context, id string) (string, error) {
	query := `SELECT name FROM users WHERE id = ?`
	row := database.QueryRowContext(ctx, query, id)
	var name string
	err := row.Scan(&name)
	return name, err
}

type homePageData struct {
	Greeting string
}

func renderHomePage(templates *template.Template, data HomePageData) []byte {
	var buffer bytes.Buffer
	err := templates.ExecuteTemplate(&buffer, "home", data)
	assert.Ok(err, "failed to render homepage")
	return buffer.Bytes()
}
