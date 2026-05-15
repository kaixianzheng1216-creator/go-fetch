package server

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/store"
)

const (
	sessionCookieName = "go_fetch_session"
	sessionUserIDKey  = "user_id"
)

func newSessionManager(dataStore *store.Store) *scs.SessionManager {
	sessions := scs.New()

	sessions.Store = pgxstore.NewWithConfig(dataStore.Pool(), pgxstore.Config{
		TableName:       "app_sessions",
		CleanUpInterval: 10 * time.Minute,
	})

	sessions.Cookie.Name = sessionCookieName
	sessions.Cookie.Secure = true
	sessions.Cookie.HttpOnly = true
	sessions.Cookie.SameSite = http.SameSiteLaxMode

	sessions.Lifetime = 24 * time.Hour

	return sessions
}
