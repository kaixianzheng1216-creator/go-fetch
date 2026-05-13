package config

import "testing"

func TestLoadRequiredConfig(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("ADMIN_USERNAME", "admin")
	t.Setenv("ADMIN_PASSWORD", "change-me")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	if cfg.DatabaseURL != "postgres://example" {
		t.Fatalf("DatabaseURL = %q", cfg.DatabaseURL)
	}

	if cfg.ListenAddr != ":8080" {
		t.Fatalf("ListenAddr = %q", cfg.ListenAddr)
	}

	if cfg.AdminUsername != "admin" {
		t.Fatalf("AdminUsername = %q", cfg.AdminUsername)
	}

	if cfg.AdminPassword != "change-me" {
		t.Fatalf("AdminPassword = %q", cfg.AdminPassword)
	}
}

func TestLoadOverridesListenAddr(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("LISTEN_ADDR", ":3000")
	t.Setenv("ADMIN_USERNAME", "admin")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	if cfg.ListenAddr != ":3000" {
		t.Fatalf("ListenAddr = %q", cfg.ListenAddr)
	}
}

func TestLoadRequiresDatabaseURL(t *testing.T) {
	_, err := Load()
	if err == nil {
		t.Fatal("expected DATABASE_URL error")
	}
}

func TestLoadRequiresAdminUsername(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://example")

	_, err := Load()
	if err == nil {
		t.Fatal("expected ADMIN_USERNAME error")
	}
}
