package user

import (
	"net/http"

	"github.com/skaisanlahti/message-board/internal/pkg/htmx"
	"github.com/skaisanlahti/message-board/internal/pkg/session"
)

type LogoutHandler struct {
	sessionManager *session.Manager
}

func NewLogoutHandler(sessionManager *session.Manager) *LogoutHandler {
	return &LogoutHandler{
		sessionManager,
	}
}

func (handler *LogoutHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	handler.sessionManager.ClearSession(response, request)
	htmx.Redirect(response, "/logout")
}
