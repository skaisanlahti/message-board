package user

import (
	"net/http"

	"github.com/skaisanlahti/message-board/internal/app/web"
)

type LogoutPageHandler struct {
	webService *web.Service
}

func NewLogoutPageHandler(
	webService *web.Service,
) *LogoutPageHandler {
	return &LogoutPageHandler{
		webService,
	}
}

func (handler *LogoutPageHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	handler.webService.Render(ctx, response, "logout_page", nil)
}
