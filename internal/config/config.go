package config

import (
	"encoding/json"
	"os"
)

type configuration struct {
	ServerAddress   string `json:"serverAddress"`
	DatabaseAddress string `json:"databaseAddress"`
}

func Read(filename string) (configuration, error) {
	var config configuration
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(bytes, &config)
	return config, nil
}
