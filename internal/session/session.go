package session

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	CookieName = "go_fetch_session"
	UserIDKey  = "user_id"
)

type Config struct {
	CookieSecure bool
	Lifetime     time.Duration
}

func NewManager(databasePool *pgxpool.Pool, config Config) *scs.SessionManager {
	sessionManager := scs.New()
	sessionManager.Store = pgxstore.NewWithConfig(databasePool, pgxstore.Config{
		TableName:       "app_sessions",
		CleanUpInterval: 10 * time.Minute,
	})
	sessionManager.Cookie.Name = CookieName
	sessionManager.Cookie.Secure = config.CookieSecure
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Lifetime = config.Lifetime
	return sessionManager
}
