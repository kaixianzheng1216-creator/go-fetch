package server

import (
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/httpapi/analytics"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/httpapi/auth"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/httpapi/events"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/httpapi/websites"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/middleware"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/session"
)

func (app *App) Routes() http.Handler {
	router := chi.NewRouter()

	router.Use(chimiddleware.RealIP)
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Timeout(60 * time.Second))
	router.Use(app.sessions.LoadAndSave)

	api := humachi.New(router, humaConfig())
	registerAPIRoutes(api, app)

	router.Get("/assets/*", app.handleFrontendAsset)
	router.Get("/script.js", app.handleScript)
	router.Get("/*", app.handleFrontend)

	return router
}

func registerAPIRoutes(api huma.API, app *App) {
	api.UseMiddleware(middleware.CaptureRequest(withRequest))
	authMiddleware := huma.Middlewares{middleware.RequireAuth(api, app.currentUser, withUser)}

	authHandler := auth.New(app.store, app.sessions, session.UserIDKey, userFromContext, isNotFound)
	eventsHandler := events.New(app.store, requestFromContext, isNotFound)
	websitesHandler := websites.New(app.store, userFromContext, websiteLookupError)
	analyticsHandler := analytics.New(app.store, userFromContext, websiteLookupError)

	events.Register(api, eventsHandler)
	auth.Register(api, authHandler, authMiddleware)
	websites.Register(api, websitesHandler, authMiddleware)
	analytics.Register(api, analyticsHandler, authMiddleware)
}
