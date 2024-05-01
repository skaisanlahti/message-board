package home

import (
	"net/http"

	"github.com/skaisanlahti/message-board/internal/view"
)

func SetupRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /", homePage)
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
