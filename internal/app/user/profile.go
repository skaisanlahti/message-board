package user

import (
	"net/http"

	"github.com/skaisanlahti/message-board/internal/app/web"
)

type ProfilePageHandler struct {
	webService *web.Service
}

func NewProfilePageHandler(
	webService *web.Service,
) *ProfilePageHandler {
	return &ProfilePageHandler{
		webService,
	}
}

func (handler *ProfilePageHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	handler.webService.Render(ctx, response, "profile_page", nil)
}
