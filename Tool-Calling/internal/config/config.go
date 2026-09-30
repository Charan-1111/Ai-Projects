package config

import (
	"os"
	"path/filepath"
	"sync"
	"tool-calling/internal/models"

	"github.com/bytedance/sonic"
)

type Configuration struct {
	Env    string        `json:"env"`
	Server models.Server `json:"server"`
	once   sync.Once
}

func (c *Configuration) LoadConfig() error {
	filePath := os.Getenv("CONFIG_PATH")
	if filePath == "" {
		filePath = filepath.Join(".", "config", "config.json")
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
