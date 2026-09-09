package config

import (
	"encoding/json"
	"os"
	"semantic-search/internal/models"
	"sync"
)

type Configuration struct {
	Env             string                        `json:"env"`
	Port            string                        `json:"port"`
	AvailableModels map[string]models.ModelConfig `json:"available_models"`
	DefaultModel    string                        `json:"default_model"`
	EmbeddingModels string                        `json:"embedding_models"`
	Retries         models.Retries                `json:"retries"`
	Queries         models.Queries                `json:"queries"`
	once            sync.Once
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
