package app

import (
	"database/sql"
	"net/http"
)

func NewApp(database *sql.DB) http.Handler {
	mux := http.NewServeMux()
	mapRoutes(mux, database)
	var handler http.Handler = mux
	// apply middleware
	// handler = middleware(handler)

	return handler
}
