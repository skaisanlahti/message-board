package app

import (
	"net/http"

	"github.com/skaisanlahti/message-board/internal/assert"
)

type AppRouter struct {
	mux *http.ServeMux
}

func (r *AppRouter) Handle(pattern string, handler http.HandlerFunc) {
	assert.NotNil(handler, "provided nil handler")
	r.mux.HandleFunc(pattern, handler)
}

func (r *AppRouter) Handler() http.Handler {
	assert.NotNil(r.mux, "serve mux was not initialized")
	return r.mux
}

func NewAppRouter() *AppRouter {
	mux := http.NewServeMux()
	return &AppRouter{mux}
}
