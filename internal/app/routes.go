package app

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/skaisanlahti/message-board/internal/app/home"
	"github.com/skaisanlahti/message-board/internal/app/log"
	"github.com/skaisanlahti/message-board/internal/app/user"
	"github.com/skaisanlahti/message-board/internal/app/web"
	"github.com/skaisanlahti/message-board/internal/pkg/middleware"
	"github.com/skaisanlahti/message-board/internal/pkg/password"
	"github.com/skaisanlahti/message-board/internal/pkg/session"
)

func newServer(
	logger *slog.Logger,
	database *sql.DB,
	htmlRenderer *web.HTMLRenderer,
	passwordHasher *password.Hasher,
	sessionManager *session.Manager,
) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /static/", web.ServeStaticFiles())

	private := middleware.New(
		session.RequireUser(),
	)

	mux.Handle("GET /", home.NewHomePageHandler(
		htmlRenderer,
	))

	mux.Handle("GET /register", user.NewRegisterPageHandler(
		htmlRenderer,
	))
	mux.Handle("POST /register", user.NewRegisterHandler(
		logger,
		database,
		sessionManager,
		passwordHasher,
		htmlRenderer,
	))

	mux.Handle("GET /login", user.NewLoginPageHandler(
		htmlRenderer,
	))
	mux.Handle("POST /login", user.NewLoginHandler(
		logger,
		database,
		sessionManager,
		passwordHasher,
		htmlRenderer,
	))

	mux.Handle("GET /logout", user.NewLogoutPageHandler(
		htmlRenderer,
	))
	mux.Handle("POST /logout", private(user.NewLogoutHandler(
		sessionManager,
	)))

	mux.Handle("GET /profile", private(user.NewProfilePageHandler(
		htmlRenderer,
	)))

	global := middleware.New(
		log.LogRequestInfo(logger),
		session.AddUserToContext(sessionManager),
	)

	handler := global(mux)
	return handler
}
