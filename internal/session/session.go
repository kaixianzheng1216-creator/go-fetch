package session

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/repository"
)

const (
	CookieName = "go_fetch_session"
	UserIDKey  = "user_id"
)

func NewManager(dataStore *repository.Store) *scs.SessionManager {
	sessionManager := scs.New()
	sessionManager.Store = pgxstore.NewWithConfig(dataStore.Pool(), pgxstore.Config{
		TableName:       "app_sessions",
		CleanUpInterval: 10 * time.Minute,
	})
	sessionManager.Cookie.Name = CookieName
	sessionManager.Cookie.Secure = true
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode
	sessionManager.Lifetime = 24 * time.Hour
	return sessionManager
}
