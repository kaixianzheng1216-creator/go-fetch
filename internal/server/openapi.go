package server

import (
	"encoding/json"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

func OpenAPIJSON() ([]byte, error) {
	router := chi.NewRouter()

	api := humachi.New(router, humaConfig())

	registerAPIRoutes(api, &App{})

	return json.MarshalIndent(api.OpenAPI(), "", "  ")
}

func humaConfig() huma.Config {
	humaAPIConfig := huma.DefaultConfig("go-fetch Analytics API", "0.1.0")

	humaAPIConfig.DocsPath = "/api/docs"

	humaAPIConfig.SchemasPath = ""

	humaAPIConfig.CreateHooks = nil

	humaAPIConfig.Servers = []*huma.Server{{URL: "/"}}

	humaAPIConfig.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"sessionCookie": {
			Type: "apiKey",
			In:   "cookie",
			Name: sessionCookieName,
		},
	}

	return humaAPIConfig
}
