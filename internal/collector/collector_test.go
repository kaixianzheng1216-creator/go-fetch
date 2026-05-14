package collector

import (
	"net/http/httptest"
	"testing"
	"time"

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
	if len(FlattenData(input.Data)) != 1 {
		t.Fatalf("expected flattened data")
	}
}

func TestFlattenDataKeepsValueTypes(t *testing.T) {
	data := map[string]any{
		"plan":   "pro",
		"paid":   true,
		"amount": float64(12.5),
		"items":  []any{"a", "b"},
		"since":  "2026-05-14T12:30:00Z",
	}

	items := FlattenData(data)
	byKey := make(map[string]FlatData, len(items))
	for _, item := range items {
		byKey[item.Key] = item
	}

	if byKey["plan"].DataType != domain.EventDataTypeString {
		t.Fatalf("plan DataType = %d", byKey["plan"].DataType)
	}
	if byKey["paid"].DataType != domain.EventDataTypeBoolean {
		t.Fatalf("paid DataType = %d", byKey["paid"].DataType)
	}
	if byKey["amount"].DataType != domain.EventDataTypeNumber || byKey["amount"].NumberValue == nil {
		t.Fatalf("amount = %#v", byKey["amount"])
	}
	if byKey["items"].DataType != domain.EventDataTypeArray {
		t.Fatalf("items DataType = %d", byKey["items"].DataType)
	}
	if byKey["since"].DataType != domain.EventDataTypeDate || byKey["since"].DateValue == nil {
		t.Fatalf("since = %#v", byKey["since"])
	}
}
