package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func New(middlewares ...Middleware) Middleware {
	return func(baseHandler http.Handler) http.Handler {
		handler := baseHandler
		for _, middleware := range middlewares {
			handler = middleware(handler)
		}

		return handler
	}
}
