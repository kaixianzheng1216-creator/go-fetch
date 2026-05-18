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

// CollectionRepository persists collected analytics events.
type CollectionRepository interface {
	GetWebsiteForCollection(ctx context.Context, websiteID uuid.UUID) (domain.Website, error)
	SaveEvent(ctx context.Context, event domain.EventInput) error
}

// CollectionInput contains the request data needed to collect an event.
type CollectionInput struct {
	Request *http.Request
	Type    domain.CollectionType
	Payload domain.CollectPayload
}

type Collector struct {
	repository CollectionRepository
	clock      Clock
}

func NewCollector(repository CollectionRepository) Collector {
	return Collector{repository: repository, clock: systemClock}
}

func (service Collector) Collect(ctx context.Context, input CollectionInput) error {
	_, isSupportedCollectionType := domain.ParseCollectionType(string(input.Type))
	if !isSupportedCollectionType {
		return ErrUnsupportedCollectionType
	}

	if input.Request == nil {
		return ErrMissingHTTPRequest
	}

	if isBot(input.Request.UserAgent()) {
		return nil
	}

	website, err := service.repository.GetWebsiteForCollection(ctx, input.Payload.WebsiteID)
	if err != nil {
		return err
	}

	clock := service.clock
	if clock == nil {
		clock = systemClock
	}

	return service.repository.SaveEvent(ctx, buildEventInput(input.Request, input.Payload, website, clock()))
}

func isBot(userAgentValue string) bool {
	return useragent.Parse(userAgentValue).Bot
}
