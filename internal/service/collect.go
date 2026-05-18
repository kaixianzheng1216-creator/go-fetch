package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
)

var (
	// ErrMissingClientInfo indicates a collection request was not bound to an HTTP request.
	ErrMissingClientInfo = errors.New("missing client info")
	// ErrUnsupportedEventType indicates an unknown tracking event type.
	ErrUnsupportedEventType = errors.New("unsupported event type")
)

// CollectionRepository persists collected analytics events.
type CollectionRepository interface {
	GetWebsiteForCollection(ctx context.Context, websiteID uuid.UUID) (domain.Website, error)
	SaveEvent(ctx context.Context, event domain.EventRecord) error
}

type ClientInfo struct {
	IP        string
	UserAgent string
	Country   string
	Region    string
	City      string
}

type CollectEventInput struct {
	Client ClientInfo
	Event  domain.TrackedEvent
}

type CollectionService struct {
	repository CollectionRepository
	clock      clock
}

func NewCollectionService(repository CollectionRepository) CollectionService {
	return CollectionService{repository: repository, clock: systemClock}
}

func (svc CollectionService) CollectEvent(ctx context.Context, input CollectEventInput) error {
	if input.Client.IP == "" && input.Client.UserAgent == "" {
		return ErrMissingClientInfo
	}

	eventType, isSupportedEventType := domain.NormalizeTrackedEventType(input.Event.Type, input.Event.Name)
	if !isSupportedEventType {
		return ErrUnsupportedEventType
	}
	input.Event.Type = eventType

	client := newEventClient(input.Client, input.Event.Screen)
	if client.bot {
		return nil
	}

	website, err := svc.repository.GetWebsiteForCollection(ctx, input.Event.WebsiteID)
	if err != nil {
		return err
	}

	return svc.repository.SaveEvent(ctx, buildEventRecord(client, input.Event, website, svc.now()))
}

func (svc CollectionService) now() time.Time {
	if svc.clock == nil {
		return systemClock()
	}
	return svc.clock()
}
