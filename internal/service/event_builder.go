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
	maxDistinctIDLength = 50

	laptopMaxScreenWidth = 1280
)

func buildEventRecord(clientInfo ClientInfo, payload domain.CollectPayload, website domain.Website, now time.Time) domain.EventRecord {
	client := newEventClient(clientInfo, payload.Screen)
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
		OS:             textutil.TruncateRunes(client.operatingSystem, maxOSLength),
		Device:         textutil.TruncateRunes(client.device, maxDeviceLength),
		Screen:         textutil.TruncateRunes(payload.Screen, maxScreenLength),
		Language:       textutil.TruncateRunes(payload.Language, maxLanguageLength),
		DistinctID:     identity.distinctID,
		CreatedAt:      now,
		Data:           payload.Data,
	}
}

type eventClient struct {
	ip              string
	userAgent       string
	browser         string
	operatingSystem string
	device          string
}

func newEventClient(clientInfo ClientInfo, screen string) eventClient {
	browser, operatingSystem, device := parseUserAgent(clientInfo.UserAgent, screen)
	return eventClient{
		ip:              clientInfo.IP,
		userAgent:       clientInfo.UserAgent,
		browser:         browser,
		operatingSystem: operatingSystem,
		device:          device,
	}
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

func eventTypeFor(eventName string) domain.EventType {
	if eventName != "" {
		return domain.EventTypeCustom
	}
	return domain.EventTypePageView
}

func sessionIDFor(websiteID uuid.UUID, distinctID, clientIP, userAgent string, now time.Time) uuid.UUID {
	value := websiteID.String() + "|" + visitorIdentity(distinctID, clientIP, userAgent) + "|" + now.UTC().Format(sessionWindowFormat)
	return stableUUID(value)
}

func visitIDFor(sessionID uuid.UUID, now time.Time) uuid.UUID {
	return stableUUID(sessionID.String() + "|" + strconv.FormatInt(now.Unix()/visitWindowSeconds, 10))
}

func parseUserAgent(userAgentValue, screen string) (browser, osName, device string) {
	agent := useragent.Parse(userAgentValue)

	browser = agent.Name
	if browser == "" || agent.IsUnknown() {
		browser = "Unknown"
	}

	osName = agent.OS
	if osName == "" {
		osName = "Unknown"
	}

	switch {
	case agent.Mobile:
		device = "mobile"
	case agent.Tablet:
		device = "tablet"
	default:
		device = "desktop"
		if width, _, hasHeight := strings.Cut(screen, "x"); hasHeight {
			if screenWidth, err := strconv.Atoi(width); err == nil && screenWidth <= laptopMaxScreenWidth {
				device = "laptop"
			}
		}
	}

	return browser, osName, device
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
	domainName := strings.TrimSpace(website.DomainName)
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
