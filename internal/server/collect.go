package server

import (
	"context"
	"net/http"
	"time"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/collector"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/httpapi"

	"github.com/danielgtaylor/huma/v2"
)

func registerCollectRoutes(api huma.API, app *App) {
	collectOp := operation(http.MethodPost, "/api/collect", "collect", "Collection", http.StatusBadRequest, http.StatusInternalServerError)
	collectOp.MaxBodyBytes = 256 * 1024
	collectOp.SkipValidateBody = true

	huma.Register(api, collectOp, app.collect)
}

func (a *App) collect(ctx context.Context, input *collectInput) (*jsonBody[httpapi.OK], error) {
	collectionType, ok := domain.ParseCollectionType(string(input.Body.Type))
	if !ok {
		return nil, huma.Error400BadRequest("unsupported collection type")
	}

	input.Body.Type = httpapi.CollectionType(collectionType)
	if input.Body.Payload.WebsiteID == "" || input.Body.Payload.URL == "" {
		return nil, huma.Error400BadRequest("website and url are required")
	}

	payload := httpapi.CollectPayloadToDomain(input.Body.Payload)
	if _, err := a.store.GetWebsiteForCollection(ctx, payload.WebsiteID); err != nil {
		if isStoreNotFound(err) {
			return nil, huma.Error400BadRequest("website not found")
		}

		return nil, huma.Error500InternalServerError("failed to load website")
	}

	r := requestFromContext(ctx)
	if r == nil {
		return nil, huma.Error500InternalServerError("failed to read request")
	}

	if collector.IsBot(r.UserAgent()) {
		return &jsonBody[httpapi.OK]{Body: httpapi.OK{OK: true}}, nil
	}

	if err := a.store.SaveEvent(ctx, collector.BuildEventInput(r, payload, time.Now())); err != nil {
		return nil, huma.Error500InternalServerError("failed to save event")
	}

	return &jsonBody[httpapi.OK]{Body: httpapi.OK{OK: true}}, nil
}
