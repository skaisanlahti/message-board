package app

import (
	"github.com/skaisanlahti/message-board/internal/assert"
	"github.com/skaisanlahti/message-board/internal/config"
	"github.com/skaisanlahti/message-board/internal/core"
	"github.com/skaisanlahti/message-board/internal/home"
)

func Start() {
	config := config.Read("appsettings.json")
	router := NewAppRouter()

	setupModules(router, config)

	server := NewServer(router, config)
	server.Start()
}

func setupModules(router core.Router, config core.Configuration) {
	assert.NotNil(router, "router was nil")
	assert.NotNil(config, "config was nil")

	home.Setup(router, config)
}
