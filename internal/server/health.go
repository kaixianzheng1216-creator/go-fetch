package server

import (
	"context"
	"net/http"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/httpapi"

	"github.com/danielgtaylor/huma/v2"
)

func registerHealthRoutes(api huma.API, app *App) {
	huma.Register(api, operation(http.MethodGet, "/healthz", "health", "Health"), app.health)

	huma.Register(api, operation(http.MethodGet, "/readyz", "ready", "Health", http.StatusServiceUnavailable), app.ready)
}

func (a *App) health(_ context.Context, _ *emptyInput) (*jsonBody[httpapi.OK], error) {
	return &jsonBody[httpapi.OK]{Body: httpapi.OK{OK: true}}, nil
}

func (a *App) ready(ctx context.Context, _ *emptyInput) (*jsonBody[httpapi.OK], error) {
	if a == nil || a.store == nil {
		return nil, huma.Error503ServiceUnavailable("database unavailable")
	}

	if err := a.store.Ping(ctx); err != nil {
		return nil, huma.Error503ServiceUnavailable("database unavailable")
	}

	return &jsonBody[httpapi.OK]{Body: httpapi.OK{OK: true}}, nil
}
