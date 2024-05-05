package config

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	ServerAddress       string `json:"serverAddress"`
	DatabaseAddress     string `json:"databaseAddress"`
	MigrationsDirectory string `json:"migrationsDirectory"`
}

func Read(path string) (Configuration, error) {
	var config Configuration
	bytes, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
