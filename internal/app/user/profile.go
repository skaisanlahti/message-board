package user

import (
	"net/http"

	"github.com/skaisanlahti/message-board/internal/app/web"
)

type ProfilePageHandler struct {
	htmlRenderer *web.HTMLRenderer
}

func NewProfilePageHandler(
	htmlRenderer *web.HTMLRenderer,
) *ProfilePageHandler {
	return &ProfilePageHandler{
		htmlRenderer,
	}
}

func (handler *ProfilePageHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	handler.htmlRenderer.Render(ctx, response, "profile_page", nil)
}
