package user

import (
	"net/http"

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
	handler.sessionManager.Stop(response, request)
	response.Header().Add("HX-Location", "/logout")
	response.WriteHeader(http.StatusOK)
}
