package web

import (
	"bytes"
	"context"
	"html/template"
	"log/slog"
	"net/http"
)

type HTMLRenderer struct {
	logger    *slog.Logger
	templates *template.Template
}

func NewHTMLRenderer(
	logger *slog.Logger,
	templates *template.Template,
) *HTMLRenderer {
	return &HTMLRenderer{
		logger,
		templates,
	}
}

func (renderer *HTMLRenderer) Render(ctx context.Context, response http.ResponseWriter, template string, data any) {
	var buffer bytes.Buffer
	err := renderer.templates.ExecuteTemplate(&buffer, template, data)
	if err != nil {
		renderer.logger.ErrorContext(
			ctx,
			"failed to render template",
			slog.String("template", template),
			slog.Any("err", err),
		)
		http.Error(response, "rendering failed", http.StatusInternalServerError)
		return
	}

	response.Header().Add("Content-type", "text/html; charset=utf-8")
	response.WriteHeader(http.StatusOK)
	response.Write(buffer.Bytes())
}
