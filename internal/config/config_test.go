package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	if cfg.DatabaseURL != "postgres://go_fetch:go_fetch@localhost:5432/go_fetch?sslmode=disable" {
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

func TestLoadOverrides(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("LISTEN_ADDR", ":3000")
	t.Setenv("ADMIN_USERNAME", "root")
	t.Setenv("ADMIN_PASSWORD", "secret")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	if cfg.DatabaseURL != "postgres://example" {
		t.Fatalf("DatabaseURL = %q", cfg.DatabaseURL)
	}

	if cfg.ListenAddr != ":3000" {
		t.Fatalf("ListenAddr = %q", cfg.ListenAddr)
	}

	if cfg.AdminUsername != "root" {
		t.Fatalf("AdminUsername = %q", cfg.AdminUsername)
	}

	if cfg.AdminPassword != "secret" {
		t.Fatalf("AdminPassword = %q", cfg.AdminPassword)
	}
}

func TestLoadRejectsEmptyRequiredValues(t *testing.T) {
	tests := []struct {
		name string
		key  string
	}{
		{name: "database URL", key: "DATABASE_URL"},
		{name: "listen addr", key: "LISTEN_ADDR"},
		{name: "admin username", key: "ADMIN_USERNAME"},
		{name: "admin password", key: "ADMIN_PASSWORD"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(tt.key, " ")

			_, err := Load()
			if err == nil {
				t.Fatal("expected error")
			}
		})
	}
}
