package service

import (
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	visitorhash "github.com/kaixianzheng1216-creator/go-fetch/internal/pkg/hash"
	ua "github.com/kaixianzheng1216-creator/go-fetch/internal/pkg/useragent"
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
)

func buildEventInput(request *http.Request, payload domain.CollectPayload, now time.Time) domain.EventInput {
	userAgent := request.UserAgent()
	ip := clientIP(request)
	browser, osName, device := ua.Parse(userAgent, payload.Screen)
	pageURL := parsePageURL(payload.URL, payload.WebsiteID)
	referrerURL := parseReferrerURL(payload.Referrer, pageURL)
	distinctID := truncate(payload.DistinctID, maxDistinctIDLength)
	eventType := domain.EventTypePageView
	if payload.Name != "" {
		eventType = domain.EventTypeCustom
	}

	sessionID := visitorhash.StableUUID(payload.WebsiteID + "|" + visitorhash.VisitorIdentity(distinctID, ip, userAgent) + "|" + now.UTC().Format(sessionWindowFormat))
	visitID := visitorhash.StableUUID(sessionID + "|" + strconv.FormatInt(now.Unix()/visitWindowSeconds, 10))

	return domain.EventInput{
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

func truncate(value string, max int) string {
	if max <= 0 {
		return ""
	}

	count := 0
	for index := range value {
		if count == max {
			return value[:index]
		}

		count++
	}

	return value
}
