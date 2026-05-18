package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
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
	Country   string
	Region    string
	City      string
}

// CollectEventParams contains the data needed to collect an event.
type CollectEventParams struct {
	Client  ClientInfo
	Type    domain.CollectionType
	Payload domain.CollectPayload
}

// CollectionService collects analytics events.
type CollectionService struct {
	repository CollectionRepository
	clock      Clock
}

// NewCollectionService returns an event collection service.
func NewCollectionService(repository CollectionRepository) CollectionService {
	return CollectionService{repository: repository, clock: systemClock}
}

// CollectEvent validates and persists an analytics event.
func (svc CollectionService) CollectEvent(ctx context.Context, params CollectEventParams) error {
	_, isSupportedCollectionType := domain.ParseCollectionType(string(params.Type))
	if !isSupportedCollectionType {
		return ErrUnsupportedCollectionType
	}

	if params.Client.IP == "" && params.Client.UserAgent == "" {
		return ErrMissingClientInfo
	}

	client := newEventClient(params.Client, params.Payload.Screen)
	if client.bot {
		return nil
	}

	website, err := svc.repository.GetWebsiteForCollection(ctx, params.Payload.WebsiteID)
	if err != nil {
		return err
	}

	clock := svc.clock
	if clock == nil {
		clock = systemClock
	}

	return svc.repository.SaveEvent(ctx, buildEventRecord(client, params.Payload, website, clock()))
}
