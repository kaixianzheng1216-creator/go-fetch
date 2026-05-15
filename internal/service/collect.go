package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/model"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/tracker"
)

var (
	ErrUnsupportedCollectionType = errors.New("unsupported collection type")
	ErrMissingHTTPRequest        = errors.New("missing http request")
)

type CollectStore interface {
	GetWebsiteForCollection(ctx context.Context, websiteID string) (model.Website, error)
	SaveEvent(ctx context.Context, event model.EventInput) error
}

type Collect struct {
	store CollectStore
	now   Clock
}

func NewCollect(store CollectStore) Collect {
	return Collect{store: store, now: systemClock}
}

func (service Collect) Collect(ctx context.Context, request *http.Request, collectionType string, payload model.CollectPayload) error {
	_, isSupportedCollectionType := model.ParseCollectionType(collectionType)
	if !isSupportedCollectionType {
		return ErrUnsupportedCollectionType
	}

	if _, err := service.store.GetWebsiteForCollection(ctx, payload.WebsiteID); err != nil {
		return err
	}

	if request == nil {
		return ErrMissingHTTPRequest
	}

	if tracker.IsBot(request.UserAgent()) {
		return nil
	}

	return service.store.SaveEvent(ctx, tracker.BuildEventInput(request, payload, service.now()))
}
