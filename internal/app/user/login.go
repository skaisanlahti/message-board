package user

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/skaisanlahti/message-board/internal/app/web"
	"github.com/skaisanlahti/message-board/internal/pkg/password"
	"github.com/skaisanlahti/message-board/internal/pkg/session"
)

type LoginHandler struct {
	logger         *slog.Logger
	database       *sql.DB
	sessionManager *session.Manager
	passwordHasher *password.Hasher
	htmlRenderer   *web.HTMLRenderer
}

func NewLoginHandler(
	logger *slog.Logger,
	database *sql.DB,
	sessionManager *session.Manager,
	passwordHasher *password.Hasher,
	htmlRenderer *web.HTMLRenderer,
) *LoginHandler {
	return &LoginHandler{
		logger,
		database,
		sessionManager,
		passwordHasher,
		htmlRenderer,
	}
}

func (handler *LoginHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	username := request.FormValue("username")
	password := request.FormValue("password")
	ctx := request.Context()

	user, err := handler.getUser(ctx, username)
	if err != nil {
		handler.handleError(ctx, response, username, password)
		return
	}

	ok := handler.passwordHasher.Verify(user.password, password)
	if !ok {
		handler.handleError(ctx, response, username, password)
		return
	}

	handler.sessionManager.Start(user.id, response)
	go handler.checkPasswordOptions(ctx, user, password)
	response.Header().Add("HX-Location", "/profile")
	response.WriteHeader(http.StatusOK)
}

type getUserResult struct {
	id       int
	password string
}

func (handler *LoginHandler) getUser(ctx context.Context, username string) (getUserResult, error) {
	query := `SELECT id, password FROM users WHERE name = $1`
	row := handler.database.QueryRowContext(ctx, query, username)

	var user getUserResult
	err := row.Scan(&user.id, &user.password)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (handler *LoginHandler) handleError(ctx context.Context, response http.ResponseWriter, username, password string) {
	data := loginPageData{
		Key:      time.Now().UnixMilli(),
		Username: username,
		Password: password,
		Error:    "invalid credentials",
	}
	handler.htmlRenderer.Render(ctx, response, "login_form", data)
}

func (handler *LoginHandler) checkPasswordOptions(ctx context.Context, user getUserResult, plainPassword string) {
	ok := handler.passwordHasher.CompareOptions(user.password)
	if ok {
		return
	}

	hashedPassword := handler.passwordHasher.Hash(plainPassword)
	query := `UPDATE users SET password = $1 WHERE id = $2`
	_, err := handler.database.ExecContext(ctx, query, hashedPassword, user.id)
	if err != nil {
		handler.logger.ErrorContext(
			ctx,
			"failed to rehash password",
			slog.Int("userID", user.id),
			slog.Any("err", err),
		)
	}
}
