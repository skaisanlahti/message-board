package app

import (
	"log/slog"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

func newMiddleware(middlewares ...Middleware) Middleware {
	return func(baseHandler http.Handler) http.Handler {
		handler := baseHandler
		for _, middleware := range middlewares {
			handler = middleware(handler)
		}

		return handler
	}
}

func logRequest(handler http.Handler, logger *slog.Logger) http.Handler {
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
