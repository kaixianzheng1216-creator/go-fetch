package tracker

import (
	"net/http"
	"testing"
	"time"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/model"
)

func TestBuildEventInput(testRunner *testing.T) {
	now := time.Date(2026, 5, 15, 9, 30, 0, 0, time.UTC)
	request, err := http.NewRequest(http.MethodPost, "/api/collect", nil)
	if err != nil {
		testRunner.Fatalf("NewRequest() error = %v", err)
	}
	request.RemoteAddr = "203.0.113.10:1234"
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/120.0 Safari/537.36")

	payload := model.CollectPayload{
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

	actual := BuildEventInput(request, payload, now)

	if actual.WebsiteID != payload.WebsiteID {
		testRunner.Fatalf("WebsiteID = %q, want %q", actual.WebsiteID, payload.WebsiteID)
	}
	if actual.EventType != model.EventTypeCustom {
		testRunner.Fatalf("EventType = %d, want %d", actual.EventType, model.EventTypeCustom)
	}
	if actual.EventName != "signup" {
		testRunner.Fatalf("EventName = %q, want signup", actual.EventName)
	}
	if actual.URLPath != "/docs#intro" {
		testRunner.Fatalf("URLPath = %q, want /docs#intro", actual.URLPath)
	}
	if actual.URLQuery != "utm_source=newsletter&utm_medium=email" {
		testRunner.Fatalf("URLQuery = %q", actual.URLQuery)
	}
	if actual.ReferrerPath != "/start" {
		testRunner.Fatalf("ReferrerPath = %q, want /start", actual.ReferrerPath)
	}
	if actual.ReferrerQuery != "q=1" {
		testRunner.Fatalf("ReferrerQuery = %q, want q=1", actual.ReferrerQuery)
	}
	if actual.ReferrerDomain != "ref.example.com" {
		testRunner.Fatalf("ReferrerDomain = %q, want ref.example.com", actual.ReferrerDomain)
	}
	if actual.Hostname != "www.example.com" {
		testRunner.Fatalf("Hostname = %q, want www.example.com", actual.Hostname)
	}
	if actual.UTMSource != "newsletter" || actual.UTMMedium != "email" {
		testRunner.Fatalf("UTM fields = (%q, %q), want (newsletter, email)", actual.UTMSource, actual.UTMMedium)
	}
	if actual.Device != "laptop" {
		testRunner.Fatalf("Device = %q, want laptop", actual.Device)
	}
	if actual.DistinctID != "user-1" {
		testRunner.Fatalf("DistinctID = %q, want user-1", actual.DistinctID)
	}
	if !actual.CreatedAt.Equal(now) {
		testRunner.Fatalf("CreatedAt = %s, want %s", actual.CreatedAt, now)
	}
	if actual.SessionID == "" || actual.VisitID == "" {
		testRunner.Fatalf("SessionID and VisitID must be populated: %#v", actual)
	}

	repeated := BuildEventInput(request, payload, now)
	if actual.SessionID != repeated.SessionID || actual.VisitID != repeated.VisitID {
		testRunner.Fatalf("stable ids changed: actual (%s, %s), repeated (%s, %s)", actual.SessionID, actual.VisitID, repeated.SessionID, repeated.VisitID)
	}
}

func TestBuildEventInputDefaultsPageview(testRunner *testing.T) {
	now := time.Date(2026, 5, 15, 9, 30, 0, 0, time.UTC)
	request, err := http.NewRequest(http.MethodPost, "/api/collect", nil)
	if err != nil {
		testRunner.Fatalf("NewRequest() error = %v", err)
	}
	request.RemoteAddr = "203.0.113.10"

	actual := BuildEventInput(request, model.CollectPayload{
		WebsiteID: "8a7e7a10-7b51-43ef-9e85-874df7dd5f8b",
		URL:       "not a valid url",
	}, now)

	if actual.EventType != model.EventTypePageView {
		testRunner.Fatalf("EventType = %d, want %d", actual.EventType, model.EventTypePageView)
	}
	if actual.URLPath != "/not%20a%20valid%20url" {
		testRunner.Fatalf("URLPath = %q, want escaped relative path", actual.URLPath)
	}
	if actual.Device != "desktop" {
		testRunner.Fatalf("Device = %q, want desktop", actual.Device)
	}
}

func TestTruncate(testRunner *testing.T) {
	tests := []struct {
		name     string
		value    string
		max      int
		expected string
	}{
		{name: "keeps shorter value", value: "abc", max: 5, expected: "abc"},
		{name: "truncates", value: "abcdef", max: 3, expected: "abc"},
		{name: "empty for non-positive max", value: "abcdef", max: 0, expected: ""},
	}

	for _, testCase := range tests {
		testRunner.Run(testCase.name, func(testRunner *testing.T) {
			if actual := truncate(testCase.value, testCase.max); actual != testCase.expected {
				testRunner.Fatalf("truncate(%q, %d) = %q, want %q", testCase.value, testCase.max, actual, testCase.expected)
			}
		})
	}
}
