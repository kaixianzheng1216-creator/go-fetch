package service

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/textutil"
	"github.com/mileusna/useragent"
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

func buildEventRecord(client eventClient, payload domain.CollectPayload, website domain.Website, now time.Time) domain.EventRecord {
	eventURLs := newEventURLs(payload, website)
	urlFields := newEventURLFields(eventURLs)
	utmFields := newUTMFields(eventURLs.page.Query())
	identity := newEventIdentity(payload.WebsiteID, payload.DistinctID, client, now)

	return domain.EventRecord{
		WebsiteID:      payload.WebsiteID,
		SessionID:      identity.sessionID,
		VisitID:        identity.visitID,
		EventType:      eventTypeFor(payload.Name),
		EventName:      textutil.TruncateRunes(payload.Name, maxEventNameLength),
		URLPath:        urlFields.path,
		URLQuery:       urlFields.query,
		ReferrerPath:   urlFields.referrerPath,
		ReferrerQuery:  urlFields.referrerQuery,
		ReferrerDomain: urlFields.referrerDomain,
		PageTitle:      textutil.TruncateRunes(payload.Title, maxPageTitleLength),
		Hostname:       urlFields.hostname,
		UTMSource:      utmFields.source,
		UTMMedium:      utmFields.medium,
		UTMCampaign:    utmFields.campaign,
		UTMContent:     utmFields.content,
		UTMTerm:        utmFields.term,
		Browser:        textutil.TruncateRunes(client.browser, maxBrowserLength),
		OS:             textutil.TruncateRunes(client.os, maxOSLength),
		Device:         textutil.TruncateRunes(client.device, maxDeviceLength),
		Screen:         textutil.TruncateRunes(payload.Screen, maxScreenLength),
		Language:       textutil.TruncateRunes(payload.Language, maxLanguageLength),
		Country:        textutil.TruncateRunes(client.country, maxCountryLength),
		Region:         textutil.TruncateRunes(client.region, maxRegionLength),
		City:           textutil.TruncateRunes(client.city, maxCityLength),
		DistinctID:     identity.distinctID,
		CreatedAt:      now,
		Data:           payload.Data,
	}
}

func eventTypeFor(eventName string) domain.EventType {
	if eventName != "" {
		return domain.EventTypeCustom
	}
	return domain.EventTypePageView
}

type eventClient struct {
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

func newEventClient(clientInfo ClientInfo, screen string) eventClient {
	agent := useragent.Parse(clientInfo.UserAgent)
	return eventClient{
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

type eventURLs struct {
	page     *url.URL
	referrer *url.URL
}

func newEventURLs(payload domain.CollectPayload, website domain.Website) eventURLs {
	pageURL := parsePageURL(payload.URL, websiteFallbackHost(website))
	return eventURLs{
		page:     pageURL,
		referrer: parseReferrerURL(payload.Referrer, pageURL),
	}
}

type eventURLFields struct {
	path           string
	query          string
	referrerPath   string
	referrerQuery  string
	referrerDomain string
	hostname       string
}

func newEventURLFields(eventURLs eventURLs) eventURLFields {
	return eventURLFields{
		path:           textutil.TruncateRunes(pathWithHash(eventURLs.page), maxURLPartLength),
		query:          textutil.TruncateRunes(eventURLs.page.RawQuery, maxURLPartLength),
		referrerPath:   textutil.TruncateRunes(eventURLs.referrer.Path, maxURLPartLength),
		referrerQuery:  textutil.TruncateRunes(eventURLs.referrer.RawQuery, maxURLPartLength),
		referrerDomain: textutil.TruncateRunes(trimWWW(eventURLs.referrer.Hostname()), maxURLPartLength),
		hostname:       textutil.TruncateRunes(eventURLs.page.Hostname(), maxHostnameLength),
	}
}

func parsePageURL(rawURL, fallbackHost string) *url.URL {
	base := &url.URL{Scheme: defaultURLScheme, Host: fallbackHost}
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
	siteDomain := strings.TrimSpace(website.Domain)
	if siteDomain == "" {
		return website.ID.String()
	}

	if !strings.Contains(siteDomain, "://") && !strings.HasPrefix(siteDomain, "//") {
		siteDomain = "//" + siteDomain
	}

	if parsedURL, err := url.Parse(siteDomain); err == nil && parsedURL.Host != "" {
		return parsedURL.Host
	}

	return website.ID.String()
}

func parseReferrerURL(rawURL string, pageURL *url.URL) *url.URL {
	if rawURL == "" {
		return &url.URL{}
	}

	base := &url.URL{Scheme: defaultURLScheme, Host: pageURL.Host}
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

type utmFields struct {
	source   string
	medium   string
	campaign string
	content  string
	term     string
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

type eventIdentity struct {
	distinctID string
	sessionID  uuid.UUID
	visitID    uuid.UUID
}

func newEventIdentity(websiteID uuid.UUID, distinctID string, client eventClient, now time.Time) eventIdentity {
	distinctID = textutil.TruncateRunes(distinctID, maxDistinctIDLength)
	sessionID := sessionIDFor(websiteID, distinctID, client.ip, client.userAgent, now)
	return eventIdentity{
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
