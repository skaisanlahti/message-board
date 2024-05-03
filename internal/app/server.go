package app

import (
	"log"
	"net/http"

	"github.com/skaisanlahti/message-board/internal/assert"
	"github.com/skaisanlahti/message-board/internal/core"
)

type AppServer struct {
	server *http.Server
}

func (s *AppServer) Start() {
	log.Printf("web server started %s", s.server.Addr)
	err := s.server.ListenAndServe()
	assert.Ok(err, "failed to start server")
}

func NewServer(router core.Router, config core.Configuration) *AppServer {
	assert.NotNil(router, "router was nil")
	assert.NotNil(config, "config was nil")

	server := &http.Server{
		Addr:    config.ServerAddress(),
		Handler: router.Handler(),
	}

	return &AppServer{server}
}
