package app

import (
	"database/sql"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/skaisanlahti/message-board/internal/app/command"
	"github.com/skaisanlahti/message-board/internal/app/query"
	"github.com/skaisanlahti/message-board/internal/app/web"
)

func newRouteHandler(templates *template.Template, _ *sql.DB, logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /static/", web.ServeStaticFiles())

	mux.Handle("GET /", query.HomePage(templates, logger))
	mux.Handle("GET /sign-up", query.SignUpPage())
	mux.Handle("GET /sign-in", query.SignInPage())
	mux.Handle("GET /sign-out", query.SignOutPage())
	mux.Handle("GET /profile", query.SignOutPage())

	mux.Handle("POST /sign-up", command.SignUp())
	mux.Handle("POST /sign-in", command.SignIn())
	mux.Handle("DELETE /sign-out", command.SignOut())

	middleware := newMiddleware(
		logRequest(logger),
	)

	handler := middleware(mux)
	return handler
}
