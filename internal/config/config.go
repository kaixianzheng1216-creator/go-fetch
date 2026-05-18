package config

import "time"

type Config struct {
	DatabaseURL               string        `env:"DATABASE_URL" envDefault:"postgres://go_fetch:go_fetch@localhost:5432/go_fetch?sslmode=disable"`
	ListenAddr                string        `env:"LISTEN_ADDR" envDefault:":8080"`
	AdminUsername             string        `env:"ADMIN_USERNAME" envDefault:"admin"`
	AdminPassword             string        `env:"ADMIN_PASSWORD"`
	HTTPReadTimeout           time.Duration `env:"HTTP_READ_TIMEOUT" envDefault:"5s"`
	HTTPWriteTimeout          time.Duration `env:"HTTP_WRITE_TIMEOUT" envDefault:"10s"`
	HTTPIdleTimeout           time.Duration `env:"HTTP_IDLE_TIMEOUT" envDefault:"120s"`
	HTTPRequestTimeout        time.Duration `env:"HTTP_REQUEST_TIMEOUT" envDefault:"60s"`
	HTTPShutdownTimeout       time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" envDefault:"10s"`
	SessionLifetime           time.Duration `env:"SESSION_LIFETIME" envDefault:"24h"`
	SessionCookieSecure       bool          `env:"SESSION_COOKIE_SECURE" envDefault:"true"`
	TrustProxyHeaders         bool          `env:"TRUST_PROXY_HEADERS" envDefault:"false"`
	CollectCORSAllowedOrigins []string      `env:"COLLECT_CORS_ALLOWED_ORIGINS" envDefault:"*" envSeparator:","`
}
