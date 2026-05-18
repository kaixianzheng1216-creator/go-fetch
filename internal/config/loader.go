package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

func Load() (Config, error) {
	config, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}

	if err := config.Validate(); err != nil {
		return Config{}, err
	}

	return config, nil
}

func LoadDatabaseURL() (string, error) {
	var config struct {
		DatabaseURL string `env:"DATABASE_URL"`
	}

	if err := env.Parse(&config); err != nil {
		return "", fmt.Errorf("parse database config: %w", err)
	}

	return databaseURLOrDefault(config.DatabaseURL), nil
}
