package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/util"
)

const (
	sessionWindowFormat = "2006-01"
	visitWindowSeconds  = 1800

	defaultURLScheme = "https"
	defaultURLPath   = "/"

	maxEventNameLength  = 50
	maxURLPartLength    = 500
	maxPageTitleLength  = 500
	maxHostnameLength   = 100
	maxUTMValueLength   = 255
	maxBrowserLength    = 20
	maxOSLength         = 20
	maxDeviceLength     = 20
	maxScreenLength     = 11
	maxLanguageLength   = 35
	maxCountryLength    = 2
	maxRegionLength     = 20
	maxCityLength       = 50
	maxDistinctIDLength = 50

	laptopMaxScreenWidth = 1280
)

var (
	ErrMissingClientInfo    = errors.New("missing client info")
	ErrUnsupportedEventType = errors.New("unsupported event type")
)

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

type clock func() time.Time

func NewCollectionService(repository CollectionRepository) CollectionService {
	return CollectionService{repository: repository, clock: systemClock}
}

func (svc CollectionService) CollectEvent(ctx context.Context, input CollectEventInput) error {
	if input.Client.IP == "" && input.Client.UserAgent == "" {
		return ErrMissingClientInfo
	}

	event := input.Event
	eventType, isSupportedEventType := domain.NormalizeTrackedEventType(event.Type, event.Name)
	if !isSupportedEventType {
		return ErrUnsupportedEventType
	}
	event.Type = eventType

	client := newTrackingClient(input.Client, event.Screen)
	if client.bot {
		return nil
	}

	website, err := svc.repository.GetWebsiteForCollection(ctx, event.WebsiteID)
	if err != nil {
		return err
	}

	return svc.repository.SaveEvent(ctx, buildEventRecord(client, event, website, svc.now()))
}

func (svc CollectionService) now() time.Time {
	if svc.clock == nil {
		return systemClock()
	}
	return svc.clock()
}

func systemClock() time.Time {
	return time.Now()
}

func buildEventRecord(client trackingClient, event domain.TrackedEvent, website domain.Website, now time.Time) domain.EventRecord {
	trackingURLs := newTrackingURLs(event, website)
	urlFields := newTrackingURLFields(trackingURLs)
	utmFields := newUTMFields(trackingURLs.page.Query())
	identity := newTrackingIdentity(event.WebsiteID, event.DistinctID, client, now)

	return domain.EventRecord{
		WebsiteID:      event.WebsiteID,
		SessionID:      identity.sessionID,
		VisitID:        identity.visitID,
		EventType:      event.Type.EventType(),
		EventName:      util.TruncateRunes(event.Name, maxEventNameLength),
		URLPath:        urlFields.path,
		URLQuery:       urlFields.query,
		ReferrerPath:   urlFields.referrerPath,
		ReferrerQuery:  urlFields.referrerQuery,
		ReferrerDomain: urlFields.referrerDomain,
		PageTitle:      util.TruncateRunes(event.Title, maxPageTitleLength),
		Hostname:       urlFields.hostname,
		UTMSource:      utmFields.source,
		UTMMedium:      utmFields.medium,
		UTMCampaign:    utmFields.campaign,
		UTMContent:     utmFields.content,
		UTMTerm:        utmFields.term,
		Browser:        util.TruncateRunes(client.browser, maxBrowserLength),
		OS:             util.TruncateRunes(client.os, maxOSLength),
		Device:         util.TruncateRunes(client.device, maxDeviceLength),
		Screen:         util.TruncateRunes(event.Screen, maxScreenLength),
		Language:       util.TruncateRunes(event.Language, maxLanguageLength),
		Country:        util.TruncateRunes(client.country, maxCountryLength),
		Region:         util.TruncateRunes(client.region, maxRegionLength),
		City:           util.TruncateRunes(client.city, maxCityLength),
		DistinctID:     identity.distinctID,
		CreatedAt:      now,
		Data:           event.Data,
	}
}
