package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func Init(path string) (*Config, error) {
	var cfg Config

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}

	if port := os.Getenv("TODO_PORT"); port != "" {
		if port[0] != ':' {
			port = ":" + port
		}
		cfg.Port = port
	}

	return &cfg, nil
}
