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

type RegisterHandler struct {
	logger          *slog.Logger
	database        *sql.DB
	sessionService  *session.Service
	passwordService *password.Service
	webService      *web.Service
}

func NewRegisterHandler(
	logger *slog.Logger,
	database *sql.DB,
	sessionService *session.Service,
	passwordService *password.Service,
	webService *web.Service,
) *RegisterHandler {
	return &RegisterHandler{
		logger,
		database,
		sessionService,
		passwordService,
		webService,
	}
}

func (handler *RegisterHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	username := request.FormValue("username")
	plainPassword := request.FormValue("password")
	ctx := request.Context()

	hashedPassword := handler.passwordService.Hash(plainPassword)
	userID, err := handler.createUser(ctx, username, hashedPassword)
	if err != nil {
		data := registerPageData{
			Key:      time.Now().UnixMilli(),
			Username: username,
			Password: plainPassword,
			Error:    "username already exists",
		}
		handler.webService.Render(ctx, response, "register_form", data)
		return
	}

	handler.sessionService.Start(userID, response)
	response.Header().Add("HX-Location", "/profile")
	response.WriteHeader(http.StatusOK)
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
