package user

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"time"

	"github.com/skaisanlahti/message-board/internal/app/web"
	"github.com/skaisanlahti/message-board/internal/pkg/htmx"
	"github.com/skaisanlahti/message-board/internal/pkg/password"
	"github.com/skaisanlahti/message-board/internal/pkg/session"
)

type RegisterHandler struct {
	logger         *slog.Logger
	database       *sql.DB
	sessionManager *session.Manager
	passwordHasher *password.Hasher
	htmlRenderer   *web.HTMLRenderer
}

func NewRegisterHandler(
	logger *slog.Logger,
	database *sql.DB,
	sessionManager *session.Manager,
	passwordHasher *password.Hasher,
	htmlRenderer *web.HTMLRenderer,
) *RegisterHandler {
	return &RegisterHandler{
		logger,
		database,
		sessionManager,
		passwordHasher,
		htmlRenderer,
	}
}

func (handler *RegisterHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	username := request.FormValue("username")
	plainPassword := request.FormValue("password")

	hashedPassword := handler.passwordHasher.Hash(plainPassword)
	userID, err := handler.createUser(ctx, username, hashedPassword)
	if err != nil {
		data := registerPageData{
			Key:      time.Now().UnixMilli(),
			Username: username,
			Password: plainPassword,
			Error:    "username already exists",
		}
		handler.htmlRenderer.Render(response, "register_form", data)
		return
	}

	handler.sessionManager.StartSession(userID, response)
	htmx.Redirect(response, "/profile")
}

func (handler *RegisterHandler) createUser(ctx context.Context, username, password string) (int, error) {
	query := `INSERT INTO users (name, password) VALUES($1, $2) RETURNING id`
	row := handler.database.QueryRowContext(ctx, query, username, password)

	var userID int
	err := row.Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
