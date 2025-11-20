package configs

import (
	"fmt"
	"os"
	"rag-pipeline/internal/utils"

	"gopkg.in/yaml.v3"
)

func NewConfig() (*Config, error) {
	config := Config{}
	if err := config.LoadConfig(); err != nil {
		return nil, err
	}
	return &config, nil
}

func (config *Config) LoadConfig() error {
	config_path := "./configs/config.prod.yaml"

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	// Load .env only in development/local
	if env == "development" || env == "dev" {
		config_path = "./configs/config.yaml"
	}

	raw_config, err := utils.LoadDocument(config_path)
	if err != nil {
		return fmt.Errorf("main.go|failed to load config.yaml: %w", err)
	}

	if err := yaml.Unmarshal([]byte(raw_config), &config); err != nil {
		return fmt.Errorf("main.go|failed to convert raw_config to yaml: %w", err)
	}

	return nil
}
