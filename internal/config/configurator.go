package config

import (
	"encoding/json"
	"os"

	"github.com/skaisanlahti/message-board/internal/assert"
)

func Read(filename string) *appConfiguration {
	bytes, err := os.ReadFile(filename)
	assert.Ok(err, "failed to read configuration file")
	assert.True(len(bytes) > 0, "configuration file has no content")

	var config jsonConfiguration
	err = json.Unmarshal(bytes, &config)
	assert.Ok(err, "failed to marshal json")

	return &appConfiguration{config}
}

type jsonConfiguration struct {
	ServerAddress   string `json:"serverAddress"`
	DatabaseAddress string `json:"databaseAddress"`
}

type appConfiguration struct {
	configuration jsonConfiguration
}

func (c *appConfiguration) ServerAddress() string {
	return c.configuration.ServerAddress
}
func (c *appConfiguration) DatabaseAddress() string {
	return c.configuration.DatabaseAddress
}
