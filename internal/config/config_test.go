package config

import (
	"testing"
	"time"
)

func TestLoadDefaultsProductionCookieSecure(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("APP_ENV", "production")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.CookieSecure {
		t.Fatal("CookieSecure = false")
	}
	if cfg.ReadTimeout != 10*time.Second {
		t.Fatalf("ReadTimeout = %s", cfg.ReadTimeout)
	}
	if cfg.DBMaxConns != 10 {
		t.Fatalf("DBMaxConns = %d", cfg.DBMaxConns)
	}
}

func TestLoadOverrides(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("COOKIE_SECURE", "false")
	t.Setenv("HTTP_READ_TIMEOUT", "3s")
	t.Setenv("DB_MAX_CONNS", "7")
	t.Setenv("LOGIN_RATE_LIMIT_PER_MINUTE", "4")
	t.Setenv("COLLECT_RATE_LIMIT_PER_MINUTE", "50")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.CookieSecure {
		t.Fatal("CookieSecure = true")
	}
	if cfg.ReadTimeout != 3*time.Second {
		t.Fatalf("ReadTimeout = %s", cfg.ReadTimeout)
	}
	if cfg.DBMaxConns != 7 {
		t.Fatalf("DBMaxConns = %d", cfg.DBMaxConns)
	}
	if cfg.LoginRateLimit != 4 {
		t.Fatalf("LoginRateLimit = %d", cfg.LoginRateLimit)
	}
	if cfg.CollectRateLimit != 50 {
		t.Fatalf("CollectRateLimit = %d", cfg.CollectRateLimit)
	}
}
