package collector

import (
	"encoding/json"
	"fmt"
	"math"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-fetch/internal/domain"

	"github.com/google/uuid"
	"github.com/mileusna/useragent"
)

const (
	sessionWindowFormat = "2006-01"
	visitWindowSeconds  = 1800

	defaultURLScheme = "https"
	defaultURLPath   = "/"

	laptopMaxScreenWidth = 1280

	maxEventNameLength = 50
	maxURLPartLength   = 500
	maxPageTitleLength = 500
	maxHostnameLength  = 100
	maxUTMValueLength  = 255
	maxScreenLength    = 20
	maxLanguageLength  = 35
	maxDataValueLength = 500
)

type FlatData struct {
	Key         string
	StringValue string
	NumberValue *float64
}

func BuildEventInput(r *http.Request, payload domain.CollectPayload, now time.Time) domain.EventInput {
	userAgent := r.UserAgent()
	clientIP := clientIP(r)
	browser, osName, device := parseUserAgent(userAgent, payload.Screen)
	pageURL := parseURL(payload.URL, payload)
	refURL := parseURL(payload.Referrer, payload)
	eventType := domain.EventTypePageView
	if payload.Name != "" {
		eventType = domain.EventTypeCustom
	}

	sessionID := stableUUID(payload.WebsiteID + "|" + clientIP + "|" + userAgent + "|" + now.UTC().Format(sessionWindowFormat))
	visitID := stableUUID(sessionID + "|" + strconv.FormatInt(now.Unix()/visitWindowSeconds, 10))

	return domain.EventInput{
		WebsiteID:      payload.WebsiteID,
		SessionID:      sessionID,
		VisitID:        visitID,
		EventType:      eventType,
		EventName:      truncate(payload.Name, maxEventNameLength),
		URLPath:        truncate(pathWithHash(pageURL), maxURLPartLength),
		URLQuery:       truncate(pageURL.RawQuery, maxURLPartLength),
		ReferrerPath:   truncate(refURL.Path, maxURLPartLength),
		ReferrerDomain: truncate(trimWWW(refURL.Hostname()), maxURLPartLength),
		PageTitle:      truncate(payload.Title, maxPageTitleLength),
		Hostname:       truncate(pageURL.Hostname(), maxHostnameLength),
		UTMSource:      truncate(pageURL.Query().Get("utm_source"), maxUTMValueLength),
		UTMMedium:      truncate(pageURL.Query().Get("utm_medium"), maxUTMValueLength),
		UTMCampaign:    truncate(pageURL.Query().Get("utm_campaign"), maxUTMValueLength),
		UTMContent:     truncate(pageURL.Query().Get("utm_content"), maxUTMValueLength),
		UTMTerm:        truncate(pageURL.Query().Get("utm_term"), maxUTMValueLength),
		Browser:        browser,
		OS:             osName,
		Device:         device,
		Screen:         truncate(payload.Screen, maxScreenLength),
		Language:       truncate(payload.Language, maxLanguageLength),
		CreatedAt:      now,
		Data:           payload.Data,
	}
}

func parseURL(raw string, payload domain.CollectPayload) *url.URL {
	baseHost := payload.WebsiteID
	if payload.URL != "" {
		if parsedURL, err := url.Parse(payload.URL); err == nil && parsedURL.Host != "" {
			baseHost = parsedURL.Host
		}
	}
	base := &url.URL{Scheme: defaultURLScheme, Host: baseHost}
	if raw == "" {
		return &url.URL{Scheme: base.Scheme, Host: base.Host, Path: defaultURLPath}
	}
	parsedURL, err := url.Parse(raw)
	if err != nil {
		return &url.URL{Scheme: base.Scheme, Host: base.Host, Path: defaultURLPath}
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

func IsBot(userAgent string) bool {
	return useragent.Parse(userAgent).Bot
}

func parseUserAgent(ua, screen string) (browser, osName, device string) {
	parsed := useragent.Parse(ua)
	browser = parsed.Name
	if browser == "" || parsed.IsUnknown() {
		browser = "Unknown"
	}
	osName = parsed.OS
	if osName == "" {
		osName = "Unknown"
	}

	switch {
	case parsed.Mobile:
		device = "mobile"
	case parsed.Tablet:
		device = "tablet"
	default:
		device = "desktop"
		if width, _, ok := strings.Cut(screen, "x"); ok {
			if n, err := strconv.Atoi(width); err == nil && n <= laptopMaxScreenWidth {
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

func FlattenData(data map[string]any) []FlatData {
	var result []FlatData
	var walk func(prefix string, value any)
	walk = func(prefix string, value any) {
		switch v := value.(type) {
		case map[string]any:
			for key, child := range v {
				walk(joinKey(prefix, key), child)
			}
		case []any:
			bytes, _ := json.Marshal(v)
			result = append(result, FlatData{Key: prefix, StringValue: truncate(string(bytes), maxDataValueLength)})
		case float64:
			if !math.IsNaN(v) && !math.IsInf(v, 0) {
				n := v
				result = append(result, FlatData{Key: prefix, StringValue: fmt.Sprintf("%g", v), NumberValue: &n})
			}
		case bool:
			result = append(result, FlatData{Key: prefix, StringValue: strconv.FormatBool(v)})
		case string:
			result = append(result, FlatData{Key: prefix, StringValue: truncate(v, maxDataValueLength)})
		case nil:
			result = append(result, FlatData{Key: prefix})
		default:
			result = append(result, FlatData{Key: prefix, StringValue: truncate(fmt.Sprint(v), maxDataValueLength)})
		}
	}
	for key, value := range data {
		walk(key, value)
	}
	return result
}

func joinKey(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return prefix + "." + key
}

func trimWWW(host string) string {
	return strings.TrimPrefix(host, "www.")
}

func truncate(value string, max int) string {
	if len(value) <= max {
		return value
	}
	return value[:max]
}
