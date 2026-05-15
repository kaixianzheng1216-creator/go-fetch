package server

import (
	"encoding/json"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/session"
)

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
			Name: session.CookieName,
		},
	}
	return cfg
}
