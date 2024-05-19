package user

import (
	"net/http"

	"github.com/skaisanlahti/message-board/internal/app/web"
)

type LogoutPageHandler struct {
	htmlRenderer *web.HTMLRenderer
}

func NewLogoutPageHandler(
	htmlRenderer *web.HTMLRenderer,
) *LogoutPageHandler {
	return &LogoutPageHandler{
		htmlRenderer,
	}
}

func (handler *LogoutPageHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	handler.htmlRenderer.Render(ctx, response, "logout_page", nil)
}
