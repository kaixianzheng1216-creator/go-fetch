package config

import (
	"fmt"
	"strings"

	"github.com/caarlos0/env/v11"
)

func Load() (Config, error) {
	var config Config

	if err := env.Parse(&config); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}

	if strings.TrimSpace(config.DatabaseURL) == "" {
		return Config{}, fmt.Errorf("DATABASE_URL cannot be empty")
	}
	if strings.TrimSpace(config.ListenAddr) == "" {
		return Config{}, fmt.Errorf("LISTEN_ADDR cannot be empty")
	}
	if strings.TrimSpace(config.AdminUsername) == "" {
		return Config{}, fmt.Errorf("ADMIN_USERNAME cannot be empty")
	}
	if strings.TrimSpace(config.AdminPassword) == "" {
		return Config{}, fmt.Errorf("ADMIN_PASSWORD cannot be empty")
	}

	return config, nil
}
