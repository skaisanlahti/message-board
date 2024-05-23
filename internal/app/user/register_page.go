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
	htmlRenderer *web.HTMLRenderer
}

func NewRegisterPageHandler(
	htmlRenderer *web.HTMLRenderer,
) *RegisterPageHandler {
	return &RegisterPageHandler{
		htmlRenderer,
	}
}

func (handler *RegisterPageHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	data := registerPageData{Key: time.Now().UnixMilli()}
	handler.htmlRenderer.Render(response, "register_page", data)
}
