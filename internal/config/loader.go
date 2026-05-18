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
	if config.HTTPReadTimeout <= 0 {
		return Config{}, fmt.Errorf("HTTP_READ_TIMEOUT must be positive")
	}
	if config.HTTPWriteTimeout <= 0 {
		return Config{}, fmt.Errorf("HTTP_WRITE_TIMEOUT must be positive")
	}
	if config.HTTPIdleTimeout <= 0 {
		return Config{}, fmt.Errorf("HTTP_IDLE_TIMEOUT must be positive")
	}
	if config.HTTPRequestTimeout <= 0 {
		return Config{}, fmt.Errorf("HTTP_REQUEST_TIMEOUT must be positive")
	}
	if config.HTTPShutdownTimeout <= 0 {
		return Config{}, fmt.Errorf("HTTP_SHUTDOWN_TIMEOUT must be positive")
	}
	if config.SessionLifetime <= 0 {
		return Config{}, fmt.Errorf("SESSION_LIFETIME must be positive")
	}

	config.CollectCORSAllowedOrigins = cleanStringSlice(config.CollectCORSAllowedOrigins)
	if len(config.CollectCORSAllowedOrigins) == 0 {
		return Config{}, fmt.Errorf("COLLECT_CORS_ALLOWED_ORIGINS cannot be empty")
	}

	return config, nil
}

func LoadDatabaseURL() (string, error) {
	var config struct {
		DatabaseURL string `env:"DATABASE_URL" envDefault:"postgres://go_fetch:go_fetch@localhost:5432/go_fetch?sslmode=disable"`
	}

	if err := env.Parse(&config); err != nil {
		return "", fmt.Errorf("parse database config: %w", err)
	}

	if strings.TrimSpace(config.DatabaseURL) == "" {
		return "", fmt.Errorf("DATABASE_URL cannot be empty")
	}

	return config.DatabaseURL, nil
}

func cleanStringSlice(values []string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			result = append(result, value)
		}
	}

	return result
}
