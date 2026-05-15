package server

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

func (a *App) Routes() http.Handler {
	r := chi.NewRouter()

	a.useHTTPMiddleware(r)

	api := humachi.New(r, humaConfig())
	registerAPIRoutes(api, a)

	r.Get("/assets/*", a.handleFrontendAsset)
	r.Get("/script.js", a.handleScript)
	r.Get("/*", a.handleFrontend)

	return r
}

func registerAPIRoutes(api huma.API, app *App) {
	auth := app.useAPIMiddleware(api)
	registerCollectRoutes(api, app)
	registerAuthRoutes(api, app, auth)
	registerWebsiteRoutes(api, app, auth)
	registerAnalyticsRoutes(api, app, auth)
}
