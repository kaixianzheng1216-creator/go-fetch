package config

type Config struct {
	DatabaseURL   string `env:"DATABASE_URL" envDefault:"postgres://go_fetch:go_fetch@localhost:5432/go_fetch?sslmode=disable"`
	ListenAddr    string `env:"LISTEN_ADDR" envDefault:":8080"`
	AdminUsername string `env:"ADMIN_USERNAME" envDefault:"admin"`
	AdminPassword string `env:"ADMIN_PASSWORD" envDefault:"change-me"`
}
