package config

import (
	"llm-playground/internal/models"
	"os"
	"sync"

	"github.com/bytedance/sonic"
)

type Configuration struct {
	Env             string                        `json:"env"`
	AvailableModels map[string]models.ModelConfig `json:"available_models"`
	DefaultModel    string                        `json:"default_model"`
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

		err = sonic.Unmarshal(fileBytes, &c)
		if err != nil {
			loadErr = err
			return
		}
	})

	return loadErr
}
