package collector

import (
	"net/http"
	"testing"
	"time"

	eventdomain "github.com/kaixianzheng1216-creator/go-fetch/internal/domain/event"
)

func TestBuildEventInput(t *testing.T) {
	now := time.Date(2026, 5, 15, 9, 30, 0, 0, time.UTC)
	request, err := http.NewRequest(http.MethodPost, "/api/collect", nil)
	if err != nil {
		t.Fatalf("NewRequest() error = %v", err)
	}
	request.RemoteAddr = "203.0.113.10:1234"
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/120.0 Safari/537.36")

	payload := eventdomain.CollectPayload{
		WebsiteID:  "8a7e7a10-7b51-43ef-9e85-874df7dd5f8b",
		URL:        "https://www.example.com/docs?utm_source=newsletter&utm_medium=email#intro",
		Referrer:   "https://ref.example.com/start?q=1",
		Title:      "Docs",
		Screen:     "1280x800",
		Language:   "en-US",
		DistinctID: "user-1",
		Name:       "signup",
		Data:       map[string]any{"plan": "pro"},
	}

	got := BuildEventInput(request, payload, now)

	if got.WebsiteID != payload.WebsiteID {
		t.Fatalf("WebsiteID = %q, want %q", got.WebsiteID, payload.WebsiteID)
	}
	if got.EventType != eventdomain.EventTypeCustom {
		t.Fatalf("EventType = %d, want %d", got.EventType, eventdomain.EventTypeCustom)
	}
	if got.EventName != "signup" {
		t.Fatalf("EventName = %q, want signup", got.EventName)
	}
	if got.URLPath != "/docs#intro" {
		t.Fatalf("URLPath = %q, want /docs#intro", got.URLPath)
	}
	if got.URLQuery != "utm_source=newsletter&utm_medium=email" {
		t.Fatalf("URLQuery = %q", got.URLQuery)
	}
	if got.ReferrerPath != "/start" {
		t.Fatalf("ReferrerPath = %q, want /start", got.ReferrerPath)
	}
	if got.ReferrerQuery != "q=1" {
		t.Fatalf("ReferrerQuery = %q, want q=1", got.ReferrerQuery)
	}
	if got.ReferrerDomain != "ref.example.com" {
		t.Fatalf("ReferrerDomain = %q, want ref.example.com", got.ReferrerDomain)
	}
	if got.Hostname != "www.example.com" {
		t.Fatalf("Hostname = %q, want www.example.com", got.Hostname)
	}
	if got.UTMSource != "newsletter" || got.UTMMedium != "email" {
		t.Fatalf("UTM fields = (%q, %q), want (newsletter, email)", got.UTMSource, got.UTMMedium)
	}
	if got.Device != "laptop" {
		t.Fatalf("Device = %q, want laptop", got.Device)
	}
	if got.DistinctID != "user-1" {
		t.Fatalf("DistinctID = %q, want user-1", got.DistinctID)
	}
	if !got.CreatedAt.Equal(now) {
		t.Fatalf("CreatedAt = %s, want %s", got.CreatedAt, now)
	}
	if got.SessionID == "" || got.VisitID == "" {
		t.Fatalf("SessionID and VisitID must be populated: %#v", got)
	}

	again := BuildEventInput(request, payload, now)
	if got.SessionID != again.SessionID || got.VisitID != again.VisitID {
		t.Fatalf("stable ids changed: got (%s, %s), again (%s, %s)", got.SessionID, got.VisitID, again.SessionID, again.VisitID)
	}
}

func TestBuildEventInputDefaultsPageview(t *testing.T) {
	now := time.Date(2026, 5, 15, 9, 30, 0, 0, time.UTC)
	request, err := http.NewRequest(http.MethodPost, "/api/collect", nil)
	if err != nil {
		t.Fatalf("NewRequest() error = %v", err)
	}
	request.RemoteAddr = "203.0.113.10"

	got := BuildEventInput(request, eventdomain.CollectPayload{
		WebsiteID: "8a7e7a10-7b51-43ef-9e85-874df7dd5f8b",
		URL:       "not a valid url",
	}, now)

	if got.EventType != eventdomain.EventTypePageView {
		t.Fatalf("EventType = %d, want %d", got.EventType, eventdomain.EventTypePageView)
	}
	if got.URLPath != "/not%20a%20valid%20url" {
		t.Fatalf("URLPath = %q, want escaped relative path", got.URLPath)
	}
	if got.Device != "desktop" {
		t.Fatalf("Device = %q, want desktop", got.Device)
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name  string
		value string
		max   int
		want  string
	}{
		{name: "keeps shorter value", value: "abc", max: 5, want: "abc"},
		{name: "truncates", value: "abcdef", max: 3, want: "abc"},
		{name: "empty for non-positive max", value: "abcdef", max: 0, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := truncate(tt.value, tt.max); got != tt.want {
				t.Fatalf("truncate(%q, %d) = %q, want %q", tt.value, tt.max, got, tt.want)
			}
		})
	}
}
