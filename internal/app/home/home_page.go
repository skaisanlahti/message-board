package home

import (
	"net/http"

	"github.com/skaisanlahti/message-board/internal/app/web"
)

type homePageData struct {
	Greeting string
}

type HomePageHandler struct {
	htmlRenderer *web.HTMLRenderer
}

func NewHomePageHandler(
	htmlRenderer *web.HTMLRenderer,
) *HomePageHandler {
	return &HomePageHandler{
		htmlRenderer,
	}
}

func (handler *HomePageHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	data := homePageData{
		Greeting: "Hello world",
	}

	handler.htmlRenderer.Render(response, "home_page", data)
}
