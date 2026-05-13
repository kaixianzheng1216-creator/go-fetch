package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	DatabaseURL      string        `env:"DATABASE_URL"`
	ListenAddr       string        `env:"LISTEN_ADDR" envDefault:":8080"`
	AdminUsername    string        `env:"ADMIN_USERNAME" envDefault:"admin"`
	AdminPassword    string        `env:"ADMIN_PASSWORD"`
	Environment      string        `env:"APP_ENV" envDefault:"development"`
	CookieSecure     bool          `env:"COOKIE_SECURE"`
	ReadTimeout      time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"10s"`
	WriteTimeout     time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"30s"`
	IdleTimeout      time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"60s"`
	HandlerTimeout   time.Duration `env:"HTTP_HANDLER_TIMEOUT" envDefault:"30s"`
	ShutdownTimeout  time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" envDefault:"10s"`
	DBMaxConns       int32         `env:"DB_MAX_CONNS" envDefault:"10"`
	LoginRateLimit   int           `env:"LOGIN_RATE_LIMIT_PER_MINUTE" envDefault:"10"`
	CollectRateLimit int           `env:"COLLECT_RATE_LIMIT_PER_MINUTE" envDefault:"120"`
}

func Load() (Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}

	if _, ok := os.LookupEnv("COOKIE_SECURE"); !ok && cfg.Environment == "production" {
		cfg.CookieSecure = true
	}

	if cfg.DatabaseURL == "" {
		return cfg, errors.New("DATABASE_URL is required")
	}

	if cfg.DBMaxConns < 1 {
		return cfg, fmt.Errorf("DB_MAX_CONNS must be positive")
	}

	if cfg.LoginRateLimit < 1 {
		return cfg, fmt.Errorf("LOGIN_RATE_LIMIT_PER_MINUTE must be positive")
	}

	if cfg.CollectRateLimit < 1 {
		return cfg, fmt.Errorf("COLLECT_RATE_LIMIT_PER_MINUTE must be positive")
	}

	return cfg, nil
}
