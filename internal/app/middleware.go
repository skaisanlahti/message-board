package app

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func addMiddleware(baseHandler http.Handler, middlewares ...Middleware) http.Handler {
	handler := baseHandler
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}

	return handler
}
