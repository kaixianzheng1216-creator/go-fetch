package config

import (
	"fmt"
	"strings"
	"time"
)

type Config struct {
	DatabaseURL               string        `env:"DATABASE_URL,notEmpty" envDefault:"postgres://go_fetch:go_fetch@localhost:5432/go_fetch?sslmode=disable"`
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

// Validate checks derived config constraints and normalizes slice fields.
func (config *Config) Validate() error {
	if config.HTTPReadTimeout <= 0 {
		return fmt.Errorf("HTTP_READ_TIMEOUT must be positive")
	}
	if config.HTTPWriteTimeout <= 0 {
		return fmt.Errorf("HTTP_WRITE_TIMEOUT must be positive")
	}
	if config.HTTPIdleTimeout <= 0 {
		return fmt.Errorf("HTTP_IDLE_TIMEOUT must be positive")
	}
	if config.HTTPRequestTimeout <= 0 {
		return fmt.Errorf("HTTP_REQUEST_TIMEOUT must be positive")
	}
	if config.HTTPShutdownTimeout <= 0 {
		return fmt.Errorf("HTTP_SHUTDOWN_TIMEOUT must be positive")
	}
	if config.SessionLifetime <= 0 {
		return fmt.Errorf("SESSION_LIFETIME must be positive")
	}

	config.CollectCORSAllowedOrigins = cleanStringSlice(config.CollectCORSAllowedOrigins)
	if len(config.CollectCORSAllowedOrigins) == 0 {
		return fmt.Errorf("COLLECT_CORS_ALLOWED_ORIGINS cannot be empty")
	}

	return nil
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
