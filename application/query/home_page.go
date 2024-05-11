package query

import (
	"html/template"
	"log/slog"
	"net/http"
)

type homePageData struct {
	Greeting string
}

func HomePage(templates *template.Template, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		data := homePageData{
			Greeting: "Hello world",
		}
		err := templates.ExecuteTemplate(response, "home_page", data)
		if err != nil {
			logger.Error("home page rendering failed")
		}
	})
}
