package service

import (
	"context"
	"errors"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mileusna/useragent"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/textutil"
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

type trackingClient struct {
	ip        string
	userAgent string
	browser   string
	os        string
	device    string
	bot       bool
	country   string
	region    string
	city      string
}

type trackingIdentity struct {
	distinctID string
	sessionID  uuid.UUID
	visitID    uuid.UUID
}

type trackingURLs struct {
	page     *url.URL
	referrer *url.URL
}

type trackingURLFields struct {
	path           string
	query          string
	referrerPath   string
	referrerQuery  string
	referrerDomain string
	hostname       string
}

type utmFields struct {
	source   string
	medium   string
	campaign string
	content  string
	term     string
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

	client := newTrackingClient(input.Client, input.Event.Screen)
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
		EventName:      textutil.TruncateRunes(event.Name, maxEventNameLength),
		URLPath:        urlFields.path,
		URLQuery:       urlFields.query,
		ReferrerPath:   urlFields.referrerPath,
		ReferrerQuery:  urlFields.referrerQuery,
		ReferrerDomain: urlFields.referrerDomain,
		PageTitle:      textutil.TruncateRunes(event.Title, maxPageTitleLength),
		Hostname:       urlFields.hostname,
		UTMSource:      utmFields.source,
		UTMMedium:      utmFields.medium,
		UTMCampaign:    utmFields.campaign,
		UTMContent:     utmFields.content,
		UTMTerm:        utmFields.term,
		Browser:        textutil.TruncateRunes(client.browser, maxBrowserLength),
		OS:             textutil.TruncateRunes(client.os, maxOSLength),
		Device:         textutil.TruncateRunes(client.device, maxDeviceLength),
		Screen:         textutil.TruncateRunes(event.Screen, maxScreenLength),
		Language:       textutil.TruncateRunes(event.Language, maxLanguageLength),
		Country:        textutil.TruncateRunes(client.country, maxCountryLength),
		Region:         textutil.TruncateRunes(client.region, maxRegionLength),
		City:           textutil.TruncateRunes(client.city, maxCityLength),
		DistinctID:     identity.distinctID,
		CreatedAt:      now,
		Data:           event.Data,
	}
}

func newTrackingClient(clientInfo ClientInfo, screen string) trackingClient {
	agent := useragent.Parse(clientInfo.UserAgent)
	return trackingClient{
		ip:        clientInfo.IP,
		userAgent: clientInfo.UserAgent,
		browser:   browserName(agent),
		os:        operatingSystemName(agent),
		device:    deviceType(agent, screen),
		bot:       agent.Bot,
		country:   clientInfo.Country,
		region:    clientInfo.Region,
		city:      clientInfo.City,
	}
}

func browserName(agent useragent.UserAgent) string {
	browser := agent.Name
	if browser == "" || agent.IsUnknown() {
		return "Unknown"
	}
	return browser
}

func operatingSystemName(agent useragent.UserAgent) string {
	if agent.OS == "" {
		return "Unknown"
	}
	return agent.OS
}

func deviceType(agent useragent.UserAgent, screen string) string {
	switch {
	case agent.Mobile:
		return "mobile"
	case agent.Tablet:
		return "tablet"
	}

	if width, ok := screenWidth(screen); ok && width <= laptopMaxScreenWidth {
		return "laptop"
	}

	return "desktop"
}

func screenWidth(screen string) (int, bool) {
	width, _, hasHeight := strings.Cut(screen, "x")
	if !hasHeight {
		return 0, false
	}

	value, err := strconv.Atoi(width)
	if err != nil {
		return 0, false
	}

	return value, true
}

func newTrackingIdentity(websiteID uuid.UUID, distinctID string, client trackingClient, now time.Time) trackingIdentity {
	distinctID = textutil.TruncateRunes(distinctID, maxDistinctIDLength)
	sessionID := sessionIDFor(websiteID, distinctID, client.ip, client.userAgent, now)
	return trackingIdentity{
		distinctID: distinctID,
		sessionID:  sessionID,
		visitID:    visitIDFor(sessionID, now),
	}
}

func sessionIDFor(websiteID uuid.UUID, distinctID, clientIP, userAgent string, now time.Time) uuid.UUID {
	value := websiteID.String() + "|" + visitorIdentity(distinctID, clientIP, userAgent) + "|" + now.UTC().Format(sessionWindowFormat)
	return stableUUID(value)
}

func visitIDFor(sessionID uuid.UUID, now time.Time) uuid.UUID {
	return stableUUID(sessionID.String() + "|" + strconv.FormatInt(now.Unix()/visitWindowSeconds, 10))
}

func stableUUID(value string) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(value))
}

func visitorIdentity(distinctID, clientIP, userAgent string) string {
	if distinctID != "" {
		return distinctID
	}
	return clientIP + "|" + userAgent
}

func newTrackingURLs(event domain.TrackedEvent, website domain.Website) trackingURLs {
	pageURL := parsePageURL(event.URL, websiteFallbackHost(website))
	return trackingURLs{
		page:     pageURL,
		referrer: parseReferrerURL(event.Referrer, pageURL),
	}
}

func newTrackingURLFields(trackingURLs trackingURLs) trackingURLFields {
	return trackingURLFields{
		path:           textutil.TruncateRunes(pathWithHash(trackingURLs.page), maxURLPartLength),
		query:          textutil.TruncateRunes(trackingURLs.page.RawQuery, maxURLPartLength),
		referrerPath:   textutil.TruncateRunes(trackingURLs.referrer.Path, maxURLPartLength),
		referrerQuery:  textutil.TruncateRunes(trackingURLs.referrer.RawQuery, maxURLPartLength),
		referrerDomain: textutil.TruncateRunes(trimWWW(trackingURLs.referrer.Hostname()), maxURLPartLength),
		hostname:       textutil.TruncateRunes(trackingURLs.page.Hostname(), maxHostnameLength),
	}
}

func parsePageURL(rawURL, fallbackHost string) *url.URL {
	base := url.URL{Scheme: defaultURLScheme, Host: fallbackHost}
	if rawURL == "" {
		return &url.URL{Scheme: base.Scheme, Host: base.Host, Path: defaultURLPath}
	}

	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return &url.URL{Scheme: base.Scheme, Host: base.Host, Path: defaultURLPath}
	}

	return base.ResolveReference(parsedURL)
}

func websiteFallbackHost(website domain.Website) string {
	domainName := strings.TrimSpace(website.Domain)
	if domainName == "" {
		return website.ID.String()
	}

	if !strings.Contains(domainName, "://") && !strings.HasPrefix(domainName, "//") {
		domainName = "//" + domainName
	}

	if parsedURL, err := url.Parse(domainName); err == nil && parsedURL.Host != "" {
		return parsedURL.Host
	}

	return website.ID.String()
}

func parseReferrerURL(rawURL string, pageURL *url.URL) *url.URL {
	if rawURL == "" {
		return &url.URL{}
	}

	base := url.URL{Scheme: defaultURLScheme, Host: pageURL.Host}
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return &url.URL{}
	}

	return base.ResolveReference(parsedURL)
}

func pathWithHash(pageURL *url.URL) string {
	path := pageURL.EscapedPath()
	if path == "" {
		path = defaultURLPath
	}

	if pageURL.Fragment != "" {
		path += "#" + pageURL.EscapedFragment()
	}

	return path
}

func trimWWW(host string) string {
	return strings.TrimPrefix(host, "www.")
}

func newUTMFields(values url.Values) utmFields {
	return utmFields{
		source:   textutil.TruncateRunes(values.Get("utm_source"), maxUTMValueLength),
		medium:   textutil.TruncateRunes(values.Get("utm_medium"), maxUTMValueLength),
		campaign: textutil.TruncateRunes(values.Get("utm_campaign"), maxUTMValueLength),
		content:  textutil.TruncateRunes(values.Get("utm_content"), maxUTMValueLength),
		term:     textutil.TruncateRunes(values.Get("utm_term"), maxUTMValueLength),
	}
}
