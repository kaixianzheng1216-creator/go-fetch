package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
)

const defaultDatabaseURL = "postgres://go_fetch:go_fetch@localhost:5432/go_fetch?sslmode=disable"

type Config struct {
	DatabaseURL               string        `env:"DATABASE_URL"`
	ListenAddr                string        `env:"LISTEN_ADDR,notEmpty" envDefault:":8080"`
	AdminUsername             string        `env:"ADMIN_USERNAME,notEmpty" envDefault:"admin"`
	AdminPassword             string        `env:"ADMIN_PASSWORD,required,notEmpty"`
	HTTPReadTimeout           time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"5s"`
	HTTPWriteTimeout          time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"10s"`
	HTTPIdleTimeout           time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"120s"`
	HTTPRequestTimeout        time.Duration `env:"HTTP_REQUEST_TIMEOUT" envDefault:"60s"`
	HTTPShutdownTimeout       time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" envDefault:"10s"`
	SessionLifetime           time.Duration `env:"SESSION_LIFETIME" envDefault:"24h"`
	SessionCookieSecure       bool          `env:"SESSION_COOKIE_SECURE" envDefault:"true"`
	TrustProxyHeaders         bool          `env:"TRUST_PROXY_HEADERS" envDefault:"false"`
	CollectCORSAllowedOrigins []string      `env:"COLLECT_CORS_ALLOWED_ORIGINS,notEmpty" envDefault:"*" envSeparator:","`
}

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

func (config *Config) Validate() error {
	config.normalize()

	for _, check := range []struct {
		name  string
		value time.Duration
	}{
		{name: "HTTP_READ_TIMEOUT", value: config.HTTPReadTimeout},
		{name: "HTTP_WRITE_TIMEOUT", value: config.HTTPWriteTimeout},
		{name: "HTTP_IDLE_TIMEOUT", value: config.HTTPIdleTimeout},
		{name: "HTTP_REQUEST_TIMEOUT", value: config.HTTPRequestTimeout},
		{name: "HTTP_SHUTDOWN_TIMEOUT", value: config.HTTPShutdownTimeout},
		{name: "SESSION_LIFETIME", value: config.SessionLifetime},
	} {
		if check.value <= 0 {
			return fmt.Errorf("%s must be positive", check.name)
		}
	}
	if len(config.CollectCORSAllowedOrigins) == 0 {
		return fmt.Errorf("COLLECT_CORS_ALLOWED_ORIGINS cannot be empty")
	}

	return nil
}

func (config *Config) normalize() {
	config.DatabaseURL = databaseURLOrDefault(config.DatabaseURL)
	config.CollectCORSAllowedOrigins = cleanStringSlice(config.CollectCORSAllowedOrigins)
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

func databaseURLOrDefault(value string) string {
	value = strings.TrimSpace(value)

	if value == "" {
		return defaultDatabaseURL
	}

	return value
}
