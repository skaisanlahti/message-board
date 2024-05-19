package log

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/skaisanlahti/message-board/internal/pkg/middleware"
)

func LogRequestInfo(logger *slog.Logger) middleware.Middleware {
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
