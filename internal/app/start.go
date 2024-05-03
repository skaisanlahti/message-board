package app

import (
	"github.com/skaisanlahti/message-board/internal/assert"
	"github.com/skaisanlahti/message-board/internal/config"
	"github.com/skaisanlahti/message-board/internal/core"
	"github.com/skaisanlahti/message-board/internal/home"
	"github.com/skaisanlahti/message-board/internal/view"
)

func Start() {
	config := config.Read("appsettings.json")
	storage := NewAppStorage(config)
	router := NewAppRouter()

	setupModules(router, config, storage)

	server := NewAppServer(router, config)
	server.Start()
}

func setupModules(router core.Router, config core.Configuration, storage core.Storage) {
	assert.NotNil(router, "router was nil")
	assert.NotNil(config, "config was nil")
	assert.NotNil(storage, "storage was nil")

	renderer := view.New()

	home.Setup(router, config, renderer)
}
