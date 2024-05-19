package web

import (
	"context"
	"html/template"
	"log/slog"
	"net/http"
)

type Service struct {
	logger    *slog.Logger
	templates *template.Template
}

func NewService(
	logger *slog.Logger,
	templates *template.Template,
) *Service {
	return &Service{
		logger,
		templates,
	}
}

func (service *Service) Render(ctx context.Context, response http.ResponseWriter, template string, data any) {
	err := service.templates.ExecuteTemplate(response, template, data)
	if err != nil {
		service.logger.ErrorContext(
			ctx,
			"failed to render template",
			slog.String("template", template),
			slog.Any("err", err),
		)
	}
}
