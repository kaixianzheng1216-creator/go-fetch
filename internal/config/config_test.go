package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	assertConfig(t, cfg, Config{
		DatabaseURL:   "postgres://go_fetch:go_fetch@localhost:5432/go_fetch?sslmode=disable",
		ListenAddr:    ":8080",
		AdminUsername: "admin",
		AdminPassword: "change-me",
	})
}

func TestLoadOverrides(t *testing.T) {
	want := Config{
		DatabaseURL:   "postgres://example",
		ListenAddr:    ":3000",
		AdminUsername: "root",
		AdminPassword: "secret",
	}

	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("LISTEN_ADDR", ":3000")
	t.Setenv("ADMIN_USERNAME", "root")
	t.Setenv("ADMIN_PASSWORD", "secret")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	assertConfig(t, cfg, want)
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

func assertConfig(t *testing.T, got, want Config) {
	t.Helper()

	if got != want {
		t.Fatalf("Config = %#v, want %#v", got, want)
	}
}
