package user

import (
	"net/http"
	"time"

	"github.com/skaisanlahti/message-board/internal/app/web"
)

type loginPageData struct {
	Key      int64
	Username string
	Password string
	Error    string
}

type LoginPageHandler struct {
	htmlRenderer *web.HTMLRenderer
}

func NewLoginPageHandler(
	htmlRenderer *web.HTMLRenderer,
) *LoginPageHandler {
	return &LoginPageHandler{
		htmlRenderer,
	}
}

func (handler *LoginPageHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	data := loginPageData{Key: time.Now().UnixMilli()}
	handler.htmlRenderer.Render(ctx, response, "login_page", data)
}
