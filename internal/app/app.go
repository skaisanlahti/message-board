package app

import (
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/skaisanlahti/message-board/internal/home"
)

func NewApp(database *sql.DB, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	addRoutes(
		mux,
		database,
	)

	handler := addMiddleware(
		mux,
		logRequest(logger),
	)

	return handler
}

func addRoutes(mux *http.ServeMux, database *sql.DB) {
	mux.Handle("GET /", home.HomePage(database))
}

func logRequest(logger *slog.Logger) Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			method := request.Method
			url := request.URL.String()
			start := time.Now()
			handler.ServeHTTP(response, request)
			durationMs := time.Since(start).Milliseconds()
			logger.Info(
				"request handled",
				slog.String("method", method),
				slog.String("url", url),
				slog.Int64("durationMs", durationMs),
			)
		})
	}
}
