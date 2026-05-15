package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/pkg/useragent"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/repository"
)

var (
	ErrUnsupportedCollectionType = errors.New("unsupported collection type")
	ErrMissingHTTPRequest        = errors.New("missing http request")
)

type Collect struct {
	store repository.TrackingRepository
	now   Clock
}

func NewCollect(store repository.TrackingRepository) Collect {
	return Collect{store: store, now: systemClock}
}

func (service Collect) Collect(ctx context.Context, request *http.Request, collectionType string, payload domain.CollectPayload) error {
	_, isSupportedCollectionType := domain.ParseCollectionType(collectionType)
	if !isSupportedCollectionType {
		return ErrUnsupportedCollectionType
	}

	if _, err := service.store.GetWebsiteForCollection(ctx, payload.WebsiteID); err != nil {
		return err
	}

	if request == nil {
		return ErrMissingHTTPRequest
	}

	if useragent.IsBot(request.UserAgent()) {
		return nil
	}

	return service.store.SaveEvent(ctx, buildEventInput(request, payload, service.now()))
}
