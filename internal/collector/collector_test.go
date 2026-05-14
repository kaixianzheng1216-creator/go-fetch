package collector

import (
	"net/http/httptest"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
)

func TestBuildEventInputParsesURLAndUTM(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/collect", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120.0")
	payload := domain.CollectPayload{
		WebsiteID:  "11111111-1111-1111-1111-111111111111",
		URL:        "https://example.com/docs?a=1&utm_source=newsletter&utm_medium=email#intro",
		Referrer:   "https://google.com/search?q=x",
		Title:      "Docs",
		Screen:     "1440x900",
		Language:   "en-US",
		DistinctID: "visitor-1",
	}

	input := BuildEventInput(req, payload, time.Date(2026, 5, 12, 12, 0, 0, 0, time.UTC))

	if input.URLPath != "/docs#intro" {
		t.Fatalf("URLPath = %q", input.URLPath)
	}
	if input.URLQuery != "a=1&utm_source=newsletter&utm_medium=email" {
		t.Fatalf("URLQuery = %q", input.URLQuery)
	}
	if input.UTMSource != "newsletter" || input.UTMMedium != "email" {
		t.Fatalf("UTM = %q/%q", input.UTMSource, input.UTMMedium)
	}
	if input.ReferrerDomain != "google.com" {
		t.Fatalf("ReferrerDomain = %q", input.ReferrerDomain)
	}
	if input.ReferrerQuery != "q=x" {
		t.Fatalf("ReferrerQuery = %q", input.ReferrerQuery)
	}
	if input.DistinctID != "visitor-1" {
		t.Fatalf("DistinctID = %q", input.DistinctID)
	}
	if input.EventType != domain.EventTypePageView {
		t.Fatalf("EventType = %d", input.EventType)
	}
}

func TestBuildEventInputLeavesEmptyReferrerEmpty(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/collect", nil)
	payload := domain.CollectPayload{
		WebsiteID: "11111111-1111-1111-1111-111111111111",
		URL:       "https://example.com/docs",
	}

	input := BuildEventInput(req, payload, time.Date(2026, 5, 12, 12, 0, 0, 0, time.UTC))

	if input.ReferrerPath != "" {
		t.Fatalf("ReferrerPath = %q", input.ReferrerPath)
	}

	if input.ReferrerDomain != "" {
		t.Fatalf("ReferrerDomain = %q", input.ReferrerDomain)
	}
}

func TestBuildEventInputCustomEvent(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/collect", nil)
	payload := domain.CollectPayload{
		WebsiteID: "11111111-1111-1111-1111-111111111111",
		URL:       "https://example.com/",
		Name:      "signup",
		Data:      map[string]any{"plan": "pro"},
	}

	input := BuildEventInput(req, payload, time.Now())

	if input.EventType != domain.EventTypeCustom {
		t.Fatalf("EventType = %d", input.EventType)
	}
	if input.EventName != "signup" {
		t.Fatalf("EventName = %q", input.EventName)
	}
	if len(domain.FlattenEventData(input.Data)) != 1 {
		t.Fatalf("expected flattened data")
	}
}

func TestBuildEventInputUsesDistinctIDForSession(t *testing.T) {
	req := httptest.NewRequest("POST", "/api/collect", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120.0")

	now := time.Date(2026, 5, 12, 12, 0, 0, 0, time.UTC)
	payload := domain.CollectPayload{
		WebsiteID: "11111111-1111-1111-1111-111111111111",
		URL:       "https://example.com/",
	}

	first := BuildEventInput(req, payload, now)
	second := BuildEventInput(req, payload, now)
	if first.SessionID != second.SessionID {
		t.Fatalf("same fallback identity produced different session ids")
	}

	payload.DistinctID = "visitor-1"
	identified := BuildEventInput(req, payload, now)
	if identified.SessionID == first.SessionID {
		t.Fatalf("distinct id did not affect session id")
	}
}

func TestBuildEventInputUsesRemoteAddrForFallbackSession(t *testing.T) {
	now := time.Date(2026, 5, 12, 12, 0, 0, 0, time.UTC)
	payload := domain.CollectPayload{
		WebsiteID: "11111111-1111-1111-1111-111111111111",
		URL:       "https://example.com/",
	}

	first := httptest.NewRequest("POST", "/api/collect", nil)
	first.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120.0")
	first.RemoteAddr = "203.0.113.10:1234"

	second := httptest.NewRequest("POST", "/api/collect", nil)
	second.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120.0")
	second.RemoteAddr = "203.0.113.11:1234"

	firstInput := BuildEventInput(first, payload, now)
	secondInput := BuildEventInput(second, payload, now)
	if firstInput.SessionID == secondInput.SessionID {
		t.Fatalf("remote addr did not affect fallback session id")
	}
}

func TestTruncateKeepsUTF8Boundary(t *testing.T) {
	got := truncate("你好世界", 3)

	if got != "你好世" {
		t.Fatalf("truncate = %q", got)
	}

	if !utf8.ValidString(got) {
		t.Fatalf("truncate returned invalid UTF-8: %q", got)
	}
}
