package config

import (
	"fmt"
	"strings"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	DatabaseURL   string `env:"DATABASE_URL" envDefault:"postgres://go_fetch:go_fetch@localhost:5432/go_fetch?sslmode=disable"`
	ListenAddr    string `env:"LISTEN_ADDR" envDefault:":8080"`
	AdminUsername string `env:"ADMIN_USERNAME" envDefault:"admin"`
	AdminPassword string `env:"ADMIN_PASSWORD" envDefault:"change-me"`
}

func Load() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, fmt.Errorf("解析配置失败: %w", err)
	}

	if strings.TrimSpace(cfg.DatabaseURL) == "" {
		return Config{}, fmt.Errorf("DATABASE_URL 不能为空")
	}

	if strings.TrimSpace(cfg.ListenAddr) == "" {
		return Config{}, fmt.Errorf("LISTEN_ADDR 不能为空")
	}

	if strings.TrimSpace(cfg.AdminUsername) == "" {
		return Config{}, fmt.Errorf("ADMIN_USERNAME 不能为空")
	}

	if strings.TrimSpace(cfg.AdminPassword) == "" {
		return Config{}, fmt.Errorf("ADMIN_PASSWORD 不能为空")
	}

	return cfg, nil
}
