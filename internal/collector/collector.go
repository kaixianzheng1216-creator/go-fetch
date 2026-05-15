package collector

import (
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	eventdomain "github.com/kaixianzheng1216-creator/go-fetch/internal/domain/event"

	"github.com/google/uuid"
	"github.com/mileusna/useragent"
)

const (
	sessionWindowFormat = "2006-01"
	visitWindowSeconds  = 1800

	defaultURLScheme = "https"
	defaultURLPath   = "/"

	laptopMaxScreenWidth = 1280

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
)

func BuildEventInput(r *http.Request, payload eventdomain.CollectPayload, now time.Time) eventdomain.EventInput {
	userAgent := r.UserAgent()
	ip := clientIP(r)
	browser, osName, device := parseUserAgent(userAgent, payload.Screen)
	pageURL := parsePageURL(payload.URL, payload.WebsiteID)
	referrerURL := parseReferrerURL(payload.Referrer, pageURL)
	distinctID := truncate(payload.DistinctID, maxDistinctIDLength)
	eventType := eventdomain.EventTypePageView
	if payload.Name != "" {
		eventType = eventdomain.EventTypeCustom
	}

	sessionID := stableUUID(payload.WebsiteID + "|" + visitorIdentity(distinctID, ip, userAgent) + "|" + now.UTC().Format(sessionWindowFormat))
	visitID := stableUUID(sessionID + "|" + strconv.FormatInt(now.Unix()/visitWindowSeconds, 10))

	return eventdomain.EventInput{
		WebsiteID:      payload.WebsiteID,
		SessionID:      sessionID,
		VisitID:        visitID,
		EventType:      eventType,
		EventName:      truncate(payload.Name, maxEventNameLength),
		URLPath:        truncate(pathWithHash(pageURL), maxURLPartLength),
		URLQuery:       truncate(pageURL.RawQuery, maxURLPartLength),
		ReferrerPath:   truncate(referrerURL.Path, maxURLPartLength),
		ReferrerQuery:  truncate(referrerURL.RawQuery, maxURLPartLength),
		ReferrerDomain: truncate(trimWWW(referrerURL.Hostname()), maxURLPartLength),
		PageTitle:      truncate(payload.Title, maxPageTitleLength),
		Hostname:       truncate(pageURL.Hostname(), maxHostnameLength),
		UTMSource:      truncate(pageURL.Query().Get("utm_source"), maxUTMValueLength),
		UTMMedium:      truncate(pageURL.Query().Get("utm_medium"), maxUTMValueLength),
		UTMCampaign:    truncate(pageURL.Query().Get("utm_campaign"), maxUTMValueLength),
		UTMContent:     truncate(pageURL.Query().Get("utm_content"), maxUTMValueLength),
		UTMTerm:        truncate(pageURL.Query().Get("utm_term"), maxUTMValueLength),
		Browser:        truncate(browser, maxBrowserLength),
		OS:             truncate(osName, maxOSLength),
		Device:         truncate(device, maxDeviceLength),
		Screen:         truncate(payload.Screen, maxScreenLength),
		Language:       truncate(payload.Language, maxLanguageLength),
		DistinctID:     distinctID,
		CreatedAt:      now,
		Data:           payload.Data,
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

func stableUUID(value string) string {
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(value)).String()
}

func visitorIdentity(distinctID, clientIP, userAgent string) string {
	if distinctID != "" {
		return distinctID
	}

	return clientIP + "|" + userAgent
}

func IsBot(userAgent string) bool {
	return useragent.Parse(userAgent).Bot
}

func parseUserAgent(userAgent, screen string) (browser, osName, device string) {
	agent := useragent.Parse(userAgent)

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
		if width, _, ok := strings.Cut(screen, "x"); ok {
			if screenWidth, err := strconv.Atoi(width); err == nil && screenWidth <= laptopMaxScreenWidth {
				device = "laptop"
			}
		}
	}

	return browser, osName, device
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}

	return r.RemoteAddr
}

func trimWWW(host string) string {
	return strings.TrimPrefix(host, "www.")
}

func truncate(value string, max int) string {
	if max <= 0 {
		return ""
	}

	count := 0
	for i := range value {
		if count == max {
			return value[:i]
		}

		count++
	}

	return value
}
