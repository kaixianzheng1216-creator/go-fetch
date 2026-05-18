package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/mileusna/useragent"
)

var (
	ErrUnsupportedCollectionType = errors.New("unsupported collection type")
	ErrMissingClientInfo         = errors.New("missing client info")
)

// CollectionRepository persists collected analytics events.
type CollectionRepository interface {
	GetWebsiteForCollection(ctx context.Context, websiteID uuid.UUID) (domain.Website, error)
	SaveEvent(ctx context.Context, event domain.EventRecord) error
}

// ClientInfo contains request-derived client metadata needed to collect events.
type ClientInfo struct {
	IP        string
	UserAgent string
}

// CollectionParams contains the data needed to collect an event.
type CollectionParams struct {
	Client  ClientInfo
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

func (service Collector) Collect(ctx context.Context, params CollectionParams) error {
	_, isSupportedCollectionType := domain.ParseCollectionType(string(params.Type))
	if !isSupportedCollectionType {
		return ErrUnsupportedCollectionType
	}

	if params.Client.IP == "" && params.Client.UserAgent == "" {
		return ErrMissingClientInfo
	}

	if isBot(params.Client.UserAgent) {
		return nil
	}

	website, err := service.repository.GetWebsiteForCollection(ctx, params.Payload.WebsiteID)
	if err != nil {
		return err
	}

	clock := service.clock
	if clock == nil {
		clock = systemClock
	}

	return service.repository.SaveEvent(ctx, buildEventRecord(params.Client, params.Payload, website, clock()))
}

func isBot(userAgentValue string) bool {
	return useragent.Parse(userAgentValue).Bot
}
