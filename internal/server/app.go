package server

import (
	"net/http"
	"time"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/store"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	sessionCookieName = "go_fetch_session"
	sessionUserIDKey  = "user_id"
)

type App struct {
	store    *store.Store
	sessions *scs.SessionManager
}

func New(dataStore *store.Store) *App {
	return &App{
		store:    dataStore,
		sessions: newSessionManager(dataStore),
	}
}

func newSessionManager(dataStore *store.Store) *scs.SessionManager {
	sessions := scs.New()

	sessions.Store = pgxstore.NewWithConfig(dataStore.Pool(), pgxstore.Config{
		TableName:       "app_sessions",
		CleanUpInterval: 5 * time.Minute,
	})

	sessions.Cookie.Name = sessionCookieName

	return sessions
}

func (a *App) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(a.sessions.LoadAndSave)

	api := humachi.New(r, humaConfig())
	registerAPIRoutes(api, a)

	r.Get("/assets/*", a.handleFrontendAsset)
	r.Get("/script.js", a.handleScript)
	r.Get("/*", a.handleFrontend)

	return r
}
