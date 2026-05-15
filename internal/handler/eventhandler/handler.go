package eventhandler

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/collector"
	eventdomain "github.com/kaixianzheng1216-creator/go-fetch/internal/event"
	websitedomain "github.com/kaixianzheng1216-creator/go-fetch/internal/website"
)

type Store interface {
	GetWebsiteForCollection(ctx context.Context, websiteID string) (websitedomain.Website, error)
	SaveEvent(ctx context.Context, event eventdomain.EventInput) error
}

type Handler struct {
	store              Store
	requestFromContext func(context.Context) *http.Request
	isNotFound         func(error) bool
}

func New(
	dataStore Store,
	requestFromContext func(context.Context) *http.Request,
	isNotFound func(error) bool,
) Handler {
	return Handler{
		store:              dataStore,
		requestFromContext: requestFromContext,
		isNotFound:         isNotFound,
	}
}

type collectRequest struct {
	Body CollectRequest
}

func (h Handler) Collect(ctx context.Context, request *collectRequest) (*okResponse, error) {
	collectionType, ok := eventdomain.ParseCollectionType(string(request.Body.Type))
	if !ok {
		return nil, huma.Error400BadRequest("unsupported collection type")
	}

	request.Body.Type = CollectionType(collectionType)

	payload := ToCollectPayload(request.Body.Payload)
	if _, err := h.store.GetWebsiteForCollection(ctx, payload.WebsiteID); err != nil {
		if h.isNotFound(err) {
			return nil, huma.Error400BadRequest("website not found")
		}

		return nil, huma.Error500InternalServerError("load website failed")
	}

	r := h.requestFromContext(ctx)
	if r == nil {
		return nil, huma.Error500InternalServerError("read request failed")
	}

	if collector.IsBot(r.UserAgent()) {
		return newOKResponse(), nil
	}

	if err := h.store.SaveEvent(ctx, collector.BuildEventInput(r, payload, time.Now())); err != nil {
		return nil, huma.Error500InternalServerError("save event failed")
	}

	return newOKResponse(), nil
}
