package httpapi

import (
	"encoding/json"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/session"
)

func OpenAPIJSON() ([]byte, error) {
	chiRouter := chi.NewRouter()
	humaAPI := humachi.New(chiRouter, humaConfig())
	server{config: Config{}.withDefaults()}.registerRoutes(humaAPI)

	return json.MarshalIndent(humaAPI.OpenAPI(), "", "  ")
}

func humaConfig() huma.Config {
	config := huma.DefaultConfig("go-fetch Analytics API", "0.1.0")
	config.DocsPath = "/api/docs"
	config.SchemasPath = ""
	config.CreateHooks = nil
	config.Servers = []*huma.Server{{URL: "/"}}
	config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"sessionCookie": {
			Type: "apiKey",
			In:   "cookie",
			Name: session.CookieName,
		},
	}
	return config
}

func publicOperation(method, path, operationID, summary, tag string) huma.Operation {
	return huma.Operation{
		Method:      method,
		Path:        path,
		OperationID: operationID,
		Summary:     summary,
		Tags:        []string{tag},
	}
}

func securedOperation(method, path, operationID, summary, tag string, authMiddleware huma.Middlewares) huma.Operation {
	operation := publicOperation(method, path, operationID, summary, tag)
	operation.Security = []map[string][]string{{"sessionCookie": {}}}
	operation.Middlewares = authMiddleware
	return operation
}

func enumValues(values []string) []any {
	result := make([]any, len(values))
	for i, value := range values {
		result[i] = value
	}
	return result
}
