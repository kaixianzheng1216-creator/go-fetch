package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/mileusna/useragent"
)

var (
	ErrUnsupportedCollectionType = errors.New("unsupported collection type")
	ErrMissingHTTPRequest        = errors.New("missing http request")
)

type CollectionStore interface {
	GetWebsiteForCollection(ctx context.Context, websiteID uuid.UUID) (domain.Website, error)
	SaveEvent(ctx context.Context, event domain.EventInput) error
}

type Collector struct {
	store CollectionStore
	clock Clock
}

func NewCollector(store CollectionStore) Collector {
	return Collector{store: store, clock: systemClock}
}

func (service Collector) Collect(ctx context.Context, request *http.Request, collectionType string, payload domain.CollectPayload) error {
	_, isSupportedCollectionType := domain.ParseCollectionType(collectionType)
	if !isSupportedCollectionType {
		return ErrUnsupportedCollectionType
	}

	if request == nil {
		return ErrMissingHTTPRequest
	}

	if isBot(request.UserAgent()) {
		return nil
	}

	if _, err := service.store.GetWebsiteForCollection(ctx, payload.WebsiteID); err != nil {
		return err
	}

	clock := service.clock
	if clock == nil {
		clock = systemClock
	}

	return service.store.SaveEvent(ctx, buildEventInput(request, payload, clock()))
}

func isBot(userAgentValue string) bool {
	return useragent.Parse(userAgentValue).Bot
}
