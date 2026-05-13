package config

import (
	"errors"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	DatabaseURL   string `env:"DATABASE_URL"`
	AdminPassword string `env:"ADMIN_PASSWORD"`
}

func Load() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}

	if cfg.DatabaseURL == "" {
		return cfg, errors.New("DATABASE_URL is required")
	}

	return cfg, nil
}
