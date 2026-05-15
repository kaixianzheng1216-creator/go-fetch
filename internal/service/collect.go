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

type TrackingStore interface {
	GetWebsiteForCollection(ctx context.Context, websiteID uuid.UUID) (domain.Website, error)
	SaveEvent(ctx context.Context, event domain.EventInput) error
}

type Collect struct {
	store TrackingStore
	now   Clock
}

func NewCollect(store TrackingStore) Collect {
	return Collect{store: store, now: systemClock}
}

func (service Collect) Collect(ctx context.Context, request *http.Request, collectionType string, payload domain.CollectPayload) error {
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

	return service.store.SaveEvent(ctx, buildEventInput(request, payload, service.now()))
}

func isBot(userAgentValue string) bool {
	return useragent.Parse(userAgentValue).Bot
}
