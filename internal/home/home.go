package home

import (
	"net/http"

	"github.com/skaisanlahti/message-board/internal/core"
)

const (
	viewHome         = "home"
	viewHomeGreeting = "home.greeting"
)

type homePageData struct {
	Greeting string
}

type controller struct {
	core.Renderer
}

func (c *controller) homePage(response http.ResponseWriter, request *http.Request) {
	data := homePageData{
		Greeting: "Home page",
	}
	html := c.Render(viewHome, data)
	response.WriteHeader(http.StatusOK)
	response.Write(html)
}

func newController(r core.Renderer) *controller {
	return &controller{r}
}

func Setup(router core.Router, config core.Configuration, renderer core.Renderer) {
	controller := newController(renderer)
	router.Handle("GET /", controller.homePage)
}
