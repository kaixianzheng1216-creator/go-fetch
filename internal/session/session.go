package session

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/store"
)

const (
	CookieName = "go_fetch_session"
	UserIDKey  = "user_id"
)

func NewManager(dataStore *store.Store) *scs.SessionManager {
	sessions := scs.New()

	sessions.Store = pgxstore.NewWithConfig(dataStore.Pool(), pgxstore.Config{
		TableName:       "app_sessions",
		CleanUpInterval: 10 * time.Minute,
	})

	sessions.Cookie.Name = CookieName
	sessions.Cookie.Secure = true
	sessions.Cookie.HttpOnly = true
	sessions.Cookie.SameSite = http.SameSiteLaxMode

	sessions.Lifetime = 24 * time.Hour

	return sessions
}
