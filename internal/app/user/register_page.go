package user

import (
	"net/http"
	"time"

	"github.com/skaisanlahti/message-board/internal/app/web"
)

type registerPageData struct {
	Key      int64
	Username string
	Password string
	Error    string
}

type RegisterPageHandler struct {
	webService *web.Service
}

func NewRegisterPageHandler(
	webService *web.Service,
) *RegisterPageHandler {
	return &RegisterPageHandler{
		webService,
	}
}

func (handler *RegisterPageHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	data := registerPageData{Key: time.Now().UnixMilli()}
	handler.webService.Render(ctx, response, "register_page", data)
}
