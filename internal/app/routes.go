package app

import (
	"database/sql"
	"net/http"

	"github.com/skaisanlahti/message-board/internal/home"
)

func mapRoutes(mux *http.ServeMux, database *sql.DB) {
	mux.Handle("GET /", home.HomePage(database))
}
