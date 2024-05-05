package app

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

func NewLogger(stdout io.Writer) Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			method := request.Method
			path := request.URL.String()
			start := time.Now()
			handler.ServeHTTP(response, request)
			milliseconds := time.Since(start).Milliseconds()
			fmt.Fprintf(stdout, "%s %s %dms\n", method, path, milliseconds)
		})
	}
}
