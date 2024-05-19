package user

import (
	"net/http"

	"github.com/skaisanlahti/message-board/internal/pkg/session"
)

type LogoutHandler struct {
	sessionService *session.Service
}

func NewLogoutHandler(sessionService *session.Service) *LogoutHandler {
	return &LogoutHandler{
		sessionService,
	}
}

func (handler *LogoutHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	handler.sessionService.Stop(response, request)
	response.Header().Add("HX-Location", "/logout")
	response.WriteHeader(http.StatusOK)
}
