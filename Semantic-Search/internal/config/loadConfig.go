package config

import (
	"encoding/json"
	"os"
	"sync"
)

type Configuration struct {
	Env  string `json:"env"`
	Port string `json:"port"`
	once sync.Once
}

func (c *Configuration) LoadConfig() error {
	filePath := os.Getenv("CONFIG_FILE_PATH")
	if filePath == "" {
		filePath = "config/local/config.json"
	}

	var loadErr error
	c.once.Do(func() {
		fileBytes, err := os.ReadFile(filePath)
		if err != nil {
			loadErr = err
			return
		}

		if err := json.Unmarshal(fileBytes, c); err != nil {
			loadErr = err
		}
	})

	return loadErr
}

func (c *Configuration) Address() string {
	if c.Port == "" {
		return ":8000"
	}
	return ":" + c.Port
}
