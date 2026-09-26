package config

import (
	"ask-my-documents/internal/models"
	"fmt"
	"os"
	"sync"

	"github.com/bytedance/sonic"
)

type Configuration struct {
	Env          string              `json:"evn"`
	Server       models.Server       `json:"server"`
	ExternalApis models.ExternalApis `json:"externalApis"`
	once         sync.Once
}

func (c *Configuration) LoadConfig() error {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return fmt.Errorf("Config file path is not founc")
	}

	var loadErr error

	c.once.Do(func() {
		fileBytes, err := os.ReadFile(configPath)
		if err != nil {
			loadErr = err
		}

		err = sonic.Unmarshal(fileBytes, &c)
		if err != nil {
			loadErr = err
		}
	})

	return loadErr
}
