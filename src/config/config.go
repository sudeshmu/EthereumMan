package config

import (
	"encoding/json"
	"etherman/src/logger"
	"os"
)

type Configuration struct {
	Users  string
	Groups []string
}

var configuration Configuration = Configuration{}

func init() {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&configuration)
	if err != nil {
		logger.ErrorLogger.Fatalf(err.Error())
	}
}

func Users() string {
	return configuration.Users
}
