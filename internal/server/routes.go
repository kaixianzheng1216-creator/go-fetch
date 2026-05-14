package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (a *App) Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(a.sessions.LoadAndSave)

	api := humachi.New(r, humaConfig())
	registerAPIRoutes(api, a)

	r.Get("/assets/*", a.handleFrontendAsset)
	r.Get("/script.js", a.handleScript)
	r.Get("/*", a.handleFrontend)

	return r
}

func registerAPIRoutes(api huma.API, app *App) {
	api.UseMiddleware(captureRequest)

	auth := huma.Middlewares{app.requireHumaAuth(api)}
	registerCollectRoutes(api, app)
	registerAuthRoutes(api, app, auth)
	registerWebsiteRoutes(api, app, auth)
	registerAnalyticsRoutes(api, app, auth)
}

func operation(method, path, operationID, tag string, errors ...int) huma.Operation {
	return huma.Operation{
		Method:      method,
		Path:        path,
		OperationID: operationID,
		Tags:        []string{tag},
		Errors:      errors,
	}
}

func authenticated(op huma.Operation, middlewares huma.Middlewares) huma.Operation {
	op.Security = []map[string][]string{{"sessionCookie": {}}}

	op.Middlewares = append(op.Middlewares, middlewares...)

	return op
}

func OpenAPIJSON() ([]byte, error) {
	r := chi.NewRouter()
	api := humachi.New(r, humaConfig())
	registerAPIRoutes(api, &App{})
	return json.MarshalIndent(api.OpenAPI(), "", "  ")
}

func humaConfig() huma.Config {
	cfg := huma.DefaultConfig("go-fetch Analytics API", "0.1.0")
	cfg.DocsPath = "/api/docs"
	cfg.SchemasPath = ""
	cfg.CreateHooks = nil
	cfg.Servers = []*huma.Server{{URL: "/"}}
	cfg.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"sessionCookie": {
			Type: "apiKey",
			In:   "cookie",
			Name: sessionCookieName,
		},
	}
	return cfg
}
