package app

import (
	"context"
	"net/http"
)

type sessionVerifier interface {
	SessionID(cookie *http.Cookie) (string, error)
	UserID(sessionID string) (string, error)
}

type Middleware = func(next http.Handler) http.Handler

func NewSessionGuard(verifier sessionVerifier) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			cookie, err := request.Cookie("sid")
			if err != nil {
				next.ServeHTTP(response, request)
				return
			}

			sessionID, err := verifier.SessionID(cookie)
			if err != nil {
				next.ServeHTTP(response, request)
				return
			}

			userID, err := verifier.UserID(sessionID)
			if err != nil {
				next.ServeHTTP(response, request)
				return
			}

			ctx := context.WithValue(request.Context(), "userID", userID)
			request = request.WithContext(ctx)
			next.ServeHTTP(response, request)
		})
	}
}
