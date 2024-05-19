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
	webService *web.Service,
	passwordService *password.Service,
	sessionService *session.Service,
) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /static/", web.ServeStaticFiles())

	private := middleware.New(
		session.RequireUser(),
	)

	mux.Handle("GET /", home.NewHomePageHandler(
		webService,
	))

	mux.Handle("GET /register", user.NewRegisterPageHandler(
		webService,
	))
	mux.Handle("POST /register", user.NewRegisterHandler(
		logger,
		database,
		sessionService,
		passwordService,
		webService,
	))

	mux.Handle("GET /login", user.NewLoginPageHandler(
		webService,
	))
	mux.Handle("POST /login", user.NewLoginHandler(
		logger,
		database,
		sessionService,
		passwordService,
		webService,
	))

	mux.Handle("GET /logout", user.NewLogoutPageHandler(
		webService,
	))
	mux.Handle("POST /logout", private(user.NewLogoutHandler(
		sessionService,
	)))

	mux.Handle("GET /profile", private(user.NewProfilePageHandler(
		webService,
	)))

	global := middleware.New(
		log.LogRequestInfo(logger),
		session.AddUserToContext(sessionService),
	)

	handler := global(mux)
	return handler
}
