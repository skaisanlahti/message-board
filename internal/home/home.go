package home

import (
	"net/http"

	"github.com/skaisanlahti/message-board/internal/core"
	"github.com/skaisanlahti/message-board/internal/view"
)

func Setup(router core.Router, config core.Configuration) {
	router.Handle("GET /", homePage)
}

const (
	viewHome         = "home"
	viewHomeGreeting = "home.greeting"
)

type homePageData struct {
	Greeting string
}

func homePage(response http.ResponseWriter, request *http.Request) {
	data := homePageData{
		Greeting: "Home page",
	}
	html := view.Render(viewHome, data)
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}
