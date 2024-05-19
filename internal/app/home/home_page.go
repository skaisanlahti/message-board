package home

import (
	"net/http"

	"github.com/skaisanlahti/message-board/internal/app/web"
)

type homePageData struct {
	Greeting string
}

type HomePageHandler struct {
	webService *web.Service
}

func NewHomePageHandler(
	webService *web.Service,
) *HomePageHandler {
	return &HomePageHandler{
		webService,
	}
}

func (handler *HomePageHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	data := homePageData{
		Greeting: "Hello world",
	}

	handler.webService.Render(request.Context(), response, "home_page", data)
}
