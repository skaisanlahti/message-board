package web

import (
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
	err := renderer.templates.ExecuteTemplate(response, template, data)
	if err != nil {
		renderer.logger.ErrorContext(
			ctx,
			"failed to render template",
			slog.String("template", template),
			slog.Any("err", err),
		)
	}
}
