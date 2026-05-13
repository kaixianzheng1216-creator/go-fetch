package server

import (
	"log/slog"
	"net/http"
	"time"

	"go-fetch/internal/config"
	"go-fetch/internal/store"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"
)

const (
	sessionCookieName = "go_fetch_session"
	sessionUserIDKey  = "user_id"
)

type App struct {
	cfg      config.Config
	store    *store.Store
	sessions *scs.SessionManager
}

func New(cfg config.Config, store *store.Store) (*App, error) {
	sessions := scs.New()
	sessions.Store = pgxstore.NewWithConfig(store.Pool(), pgxstore.Config{
		TableName:       "app_sessions",
		CleanUpInterval: 5 * time.Minute,
	})
	sessions.Lifetime = 7 * 24 * time.Hour
	sessions.Cookie.Name = sessionCookieName
	sessions.Cookie.HttpOnly = true
	sessions.Cookie.Path = "/"
	sessions.Cookie.SameSite = http.SameSiteLaxMode
	sessions.Cookie.Secure = cfg.CookieSecure

	return &App{cfg: cfg, store: store, sessions: sessions}, nil
}

func (a *App) Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(httplog.RequestLogger(slog.Default(), &httplog.Options{Level: slog.LevelInfo, Schema: httplog.SchemaECS}))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(a.cfg.HandlerTimeout))
	r.Use(a.secureHeaders().Handler)

	api := humachi.New(r, humaConfig())
	registerAPIRoutes(api, a)

	r.Get("/assets/*", a.handleFrontendAsset)
	r.Get("/script.js", a.handleScript)
	r.Get("/*", a.handleFrontend)
	return a.sessions.LoadAndSave(r)
}
