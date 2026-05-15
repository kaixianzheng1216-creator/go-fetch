package httpapi

import (
	"encoding/json"

	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

func OpenAPIJSON() ([]byte, error) {
	chiRouter := chi.NewRouter()
	humaAPI := humachi.New(chiRouter, humaConfig())
	server{config: Config{}.withDefaults()}.registerRoutes(humaAPI)

	return json.MarshalIndent(humaAPI.OpenAPI(), "", "  ")
}
