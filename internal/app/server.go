package app

import (
	"database/sql"
	"io"
	"net/http"
)

func NewApp(database *sql.DB, stdout io.Writer) http.Handler {
	mux := http.NewServeMux()
	logger := NewLogger(stdout)
	mapRoutes(mux, database)
	var handler http.Handler = mux
	handler = logger(handler)
	return handler
}
