package collector

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
)

const testWebsiteID = "11111111-1111-1111-1111-111111111111"

var testNow = time.Date(2026, 5, 12, 12, 0, 0, 0, time.UTC)

func eventRequest() *http.Request {
	req := httptest.NewRequest("POST", "/api/collect", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 Chrome/120.0")
	return req
}

func TestBuildEventInputParsesURLAndUTM(t *testing.T) {
	payload := domain.CollectPayload{
		WebsiteID:  testWebsiteID,
		URL:        "https://example.com/docs?a=1&utm_source=newsletter&utm_medium=email#intro",
		Referrer:   "https://google.com/search?q=x",
		Title:      "Docs",
		Screen:     "1440x900",
		Language:   "en-US",
		DistinctID: "visitor-1",
	}

	input := BuildEventInput(eventRequest(), payload, testNow)

	assertString(t, "URLPath", input.URLPath, "/docs#intro")
	assertString(t, "URLQuery", input.URLQuery, "a=1&utm_source=newsletter&utm_medium=email")
	assertString(t, "UTMSource", input.UTMSource, "newsletter")
	assertString(t, "UTMMedium", input.UTMMedium, "email")
	assertString(t, "ReferrerDomain", input.ReferrerDomain, "google.com")
	assertString(t, "ReferrerQuery", input.ReferrerQuery, "q=x")
	assertString(t, "DistinctID", input.DistinctID, "visitor-1")
	if input.EventType != domain.EventTypePageView {
		t.Fatalf("EventType = %d", input.EventType)
	}
}

func TestBuildEventInputLeavesEmptyReferrerEmpty(t *testing.T) {
	payload := domain.CollectPayload{
		WebsiteID: testWebsiteID,
		URL:       "https://example.com/docs",
	}

	input := BuildEventInput(eventRequest(), payload, testNow)

	assertString(t, "ReferrerPath", input.ReferrerPath, "")
	assertString(t, "ReferrerDomain", input.ReferrerDomain, "")
}

func TestBuildEventInputCustomEvent(t *testing.T) {
	payload := domain.CollectPayload{
		WebsiteID: testWebsiteID,
		URL:       "https://example.com/",
		Name:      "signup",
		Data:      map[string]any{"plan": "pro"},
	}

	input := BuildEventInput(eventRequest(), payload, testNow)

	if input.EventType != domain.EventTypeCustom {
		t.Fatalf("EventType = %d", input.EventType)
	}
	assertString(t, "EventName", input.EventName, "signup")

	plan, ok := input.Data["plan"].(string)
	if !ok {
		t.Fatalf("Data[plan] = %#v, want string", input.Data["plan"])
	}
	assertString(t, "Data[plan]", plan, "pro")
}

func TestBuildEventInputUsesDistinctIDForSession(t *testing.T) {
	payload := domain.CollectPayload{
		WebsiteID: testWebsiteID,
		URL:       "https://example.com/",
	}

	req := eventRequest()
	first := BuildEventInput(req, payload, testNow)
	second := BuildEventInput(req, payload, testNow)
	if first.SessionID != second.SessionID {
		t.Fatalf("same fallback identity produced different session ids")
	}

	payload.DistinctID = "visitor-1"
	identified := BuildEventInput(req, payload, testNow)
	if identified.SessionID == first.SessionID {
		t.Fatalf("distinct id did not affect session id")
	}
}

func TestBuildEventInputUsesRemoteAddrForFallbackSession(t *testing.T) {
	payload := domain.CollectPayload{
		WebsiteID: testWebsiteID,
		URL:       "https://example.com/",
	}

	first := eventRequest()
	first.RemoteAddr = "203.0.113.10:1234"

	second := eventRequest()
	second.RemoteAddr = "203.0.113.11:1234"

	firstInput := BuildEventInput(first, payload, testNow)
	secondInput := BuildEventInput(second, payload, testNow)
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

func assertString(t *testing.T, name, got, want string) {
	t.Helper()

	if got != want {
		t.Fatalf("%s = %q, want %q", name, got, want)
	}
}
