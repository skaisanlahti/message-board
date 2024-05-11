package program

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/skaisanlahti/message-board/application/command"
	"github.com/skaisanlahti/message-board/application/query"
)

func newHandler(_ *sql.DB, logger *slog.Logger) http.Handler {
	router := http.NewServeMux()
	router.Handle("GET /", query.HomePage())
	router.Handle("GET /sign-up", query.SignUpPage())
	router.Handle("GET /sign-in", query.SignInPage())
	router.Handle("GET /sign-out", query.SignOutPage())
	router.Handle("GET /profile", query.SignOutPage())

	router.Handle("POST /sign-up", command.SignUp())
	router.Handle("POST /sign-in", command.SignIn())
	router.Handle("DELETE /sign-out", command.SignOut())

	handler := logRequest(router, logger)
	return handler
}
