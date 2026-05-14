package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/store"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog/v3"
)

type App struct {
	store    *store.Store
	sessions *scs.SessionManager
}

func New(dataStore *store.Store, secureCookie bool) *App {
	sessions := scs.New()
	sessions.Lifetime = 7 * 24 * time.Hour
	sessions.Store = pgxstore.NewWithConfig(dataStore.Pool(), pgxstore.Config{
		TableName:       "app_sessions",
		CleanUpInterval: 5 * time.Minute,
	})
	sessions.Codec = scs.GobCodec{}
	sessions.Cookie.Name = "go_fetch_session"
	sessions.Cookie.Secure = secureCookie

	return &App{store: dataStore, sessions: sessions}
}

func (a *App) Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(httplog.RequestLogger(slog.Default(), &httplog.Options{Level: slog.LevelInfo, Schema: httplog.SchemaECS}))
	r.Use(middleware.Recoverer)
	r.Use(limitRequestBody(1 << 20))
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(a.sessions.LoadAndSave)
	r.Use(a.secureHeaders().Handler)

	api := humachi.New(r, humaConfig())
	registerAPIRoutes(api, a)

	r.With(middleware.SetHeader("Cache-Control", "public, max-age=31536000, immutable")).Get("/assets/*", a.handleFrontendAsset)
	r.Get("/script.js", a.handleScript)
	r.Get("/*", a.handleFrontend)
	return r
}
