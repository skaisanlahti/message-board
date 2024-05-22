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

func newRouter(
	logger *slog.Logger,
	database *sql.DB,
	htmlRenderer *web.HTMLRenderer,
	passwordHasher *password.Hasher,
	sessionManager *session.Manager,
) http.Handler {
	router := http.NewServeMux()
	router.Handle("GET /static/", web.ServeStaticFiles())

	private := middleware.New(
		session.Require(true, http.RedirectHandler("/login", http.StatusSeeOther)),
	)

	public := middleware.New(
		session.Require(false, http.NotFoundHandler()),
	)

	router.Handle("GET /", home.NewHomePageHandler(
		htmlRenderer,
	))

	router.Handle("GET /register", public(user.NewRegisterPageHandler(
		htmlRenderer,
	)))
	router.Handle("POST /register", public(user.NewRegisterHandler(
		logger,
		database,
		sessionManager,
		passwordHasher,
		htmlRenderer,
	)))

	router.Handle("GET /login", public(user.NewLoginPageHandler(
		htmlRenderer,
	)))
	router.Handle("POST /login", public(user.NewLoginHandler(
		logger,
		database,
		sessionManager,
		passwordHasher,
		htmlRenderer,
	)))

	router.Handle("GET /logout", public(user.NewLogoutPageHandler(
		htmlRenderer,
	)))
	router.Handle("POST /logout", private(user.NewLogoutHandler(
		sessionManager,
	)))

	router.Handle("GET /profile", private(user.NewProfilePageHandler(
		htmlRenderer,
	)))

	global := middleware.New(
		log.LogRequestInfo(logger),
		session.Middleware(sessionManager),
	)

	handler := global(router)
	return handler
}
