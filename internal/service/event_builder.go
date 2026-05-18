package service

import (
	"net"
	"net/http"
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

func buildEventInput(request *http.Request, payload domain.CollectPayload, now time.Time) domain.EventInput {
	client := newEventClient(request, payload.Screen)
	eventURLs := newEventURLs(payload)
	distinctID := textutil.TruncateRunes(payload.DistinctID, maxDistinctIDLength)
	sessionID := sessionIDFor(payload.WebsiteID, distinctID, client.ip, client.userAgent, now)
	visitID := visitIDFor(sessionID, now)
	pageQuery := eventURLs.page.Query()

	return domain.EventInput{
		WebsiteID:      payload.WebsiteID,
		SessionID:      sessionID,
		VisitID:        visitID,
		EventType:      eventTypeFor(payload.Name),
		EventName:      textutil.TruncateRunes(payload.Name, maxEventNameLength),
		URLPath:        textutil.TruncateRunes(pathWithHash(eventURLs.page), maxURLPartLength),
		URLQuery:       textutil.TruncateRunes(eventURLs.page.RawQuery, maxURLPartLength),
		ReferrerPath:   textutil.TruncateRunes(eventURLs.referrer.Path, maxURLPartLength),
		ReferrerQuery:  textutil.TruncateRunes(eventURLs.referrer.RawQuery, maxURLPartLength),
		ReferrerDomain: textutil.TruncateRunes(trimWWW(eventURLs.referrer.Hostname()), maxURLPartLength),
		PageTitle:      textutil.TruncateRunes(payload.Title, maxPageTitleLength),
		Hostname:       textutil.TruncateRunes(eventURLs.page.Hostname(), maxHostnameLength),
		UTMSource:      textutil.TruncateRunes(pageQuery.Get("utm_source"), maxUTMValueLength),
		UTMMedium:      textutil.TruncateRunes(pageQuery.Get("utm_medium"), maxUTMValueLength),
		UTMCampaign:    textutil.TruncateRunes(pageQuery.Get("utm_campaign"), maxUTMValueLength),
		UTMContent:     textutil.TruncateRunes(pageQuery.Get("utm_content"), maxUTMValueLength),
		UTMTerm:        textutil.TruncateRunes(pageQuery.Get("utm_term"), maxUTMValueLength),
		Browser:        textutil.TruncateRunes(client.browser, maxBrowserLength),
		OS:             textutil.TruncateRunes(client.operatingSystem, maxOSLength),
		Device:         textutil.TruncateRunes(client.device, maxDeviceLength),
		Screen:         textutil.TruncateRunes(payload.Screen, maxScreenLength),
		Language:       textutil.TruncateRunes(payload.Language, maxLanguageLength),
		DistinctID:     distinctID,
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

func newEventClient(request *http.Request, screen string) eventClient {
	userAgentValue := request.UserAgent()
	browser, operatingSystem, device := parseUserAgent(userAgentValue, screen)
	return eventClient{
		ip:              clientIP(request),
		userAgent:       userAgentValue,
		browser:         browser,
		operatingSystem: operatingSystem,
		device:          device,
	}
}

type eventURLs struct {
	page     *url.URL
	referrer *url.URL
}

func newEventURLs(payload domain.CollectPayload) eventURLs {
	pageURL := parsePageURL(payload.URL, payload.WebsiteID.String())
	return eventURLs{
		page:     pageURL,
		referrer: parseReferrerURL(payload.Referrer, pageURL),
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

func clientIP(request *http.Request) string {
	host, _, err := net.SplitHostPort(request.RemoteAddr)
	if err == nil {
		return host
	}

	return request.RemoteAddr
}

func trimWWW(host string) string {
	return strings.TrimPrefix(host, "www.")
}
