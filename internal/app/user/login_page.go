package user

import (
	"net/http"
	"time"

	"github.com/skaisanlahti/message-board/internal/app/web"
	"github.com/skaisanlahti/message-board/internal/pkg/session"
)

type loginPageData struct {
	Key      int64
	Username string
	Password string
	Error    string
}

type LoginPageHandler struct {
	webService *web.Service
}

func NewLoginPageHandler(
	webService *web.Service,
) *LoginPageHandler {
	return &LoginPageHandler{
		webService,
	}
}

func (handler *LoginPageHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	_, ok := session.GetUserFromContext(request)
	if ok {
		response.Header().Add("HX-Location", "/profile")
		response.WriteHeader(http.StatusOK)
		return
	}

	ctx := request.Context()
	data := loginPageData{Key: time.Now().UnixMilli()}
	handler.webService.Render(ctx, response, "login_page", data)
}
