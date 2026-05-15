package event

import (
	"reflect"
	"testing"
	"time"
)

func TestParseCollectionType(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  CollectionType
		ok    bool
	}{
		{name: "empty defaults to event", value: "", want: CollectionTypeEvent, ok: true},
		{name: "event", value: "event", want: CollectionTypeEvent, ok: true},
		{name: "unsupported", value: "pageview", ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ParseCollectionType(tt.value)
			if got != tt.want || ok != tt.ok {
				t.Fatalf("ParseCollectionType(%q) = (%q, %v), want (%q, %v)", tt.value, got, ok, tt.want, tt.ok)
			}
		})
	}
}

func TestParseDateUnit(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  DateUnit
	}{
		{name: "month", value: "month", want: DateUnitMonth},
		{name: "day", value: "day", want: DateUnitDay},
		{name: "hour", value: "hour", want: DefaultDateUnit},
		{name: "invalid", value: "invalid", want: DefaultDateUnit},
		{name: "empty", value: "", want: DefaultDateUnit},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseDateUnit(tt.value); got != tt.want {
				t.Fatalf("ParseDateUnit(%q) = %q, want %q", tt.value, got, tt.want)
			}
		})
	}
}

func TestDateRange(t *testing.T) {
	startMillis := time.Date(2026, 5, 15, 8, 0, 0, 0, time.UTC).UnixMilli()
	endMillis := time.Date(2026, 5, 15, 9, 0, 0, 0, time.UTC).UnixMilli()

	start, end, unit := DateRange(&startMillis, &endMillis, "day")

	if !start.Equal(time.UnixMilli(startMillis)) {
		t.Fatalf("start = %s, want %s", start, time.UnixMilli(startMillis))
	}
	if !end.Equal(time.UnixMilli(endMillis)) {
		t.Fatalf("end = %s, want %s", end, time.UnixMilli(endMillis))
	}
	if unit != DateUnitDay {
		t.Fatalf("unit = %q, want %q", unit, DateUnitDay)
	}
}

func TestFormatBucket(t *testing.T) {
	bucket := time.Date(2026, 5, 15, 9, 30, 0, 0, time.UTC)
	tests := []struct {
		unit DateUnit
		want string
	}{
		{unit: DateUnitMonth, want: "2026-05"},
		{unit: DateUnitDay, want: "05-15"},
		{unit: DateUnitHour, want: "09:30"},
	}

	for _, tt := range tests {
		t.Run(string(tt.unit), func(t *testing.T) {
			if got := FormatBucket(bucket, tt.unit); got != tt.want {
				t.Fatalf("FormatBucket(%s, %q) = %q, want %q", bucket, tt.unit, got, tt.want)
			}
		})
	}
}

func TestMetricRules(t *testing.T) {
	metric, ok := ParseMetricType("browser")
	if !ok {
		t.Fatal("ParseMetricType(browser) returned ok=false")
	}

	if metric.EventType() != EventTypePageView {
		t.Fatalf("browser EventType() = %d, want %d", metric.EventType(), EventTypePageView)
	}
	if !metric.IsSessionDimension() {
		t.Fatal("browser should be a session dimension")
	}

	eventMetric, ok := ParseMetricType("event")
	if !ok {
		t.Fatal("ParseMetricType(event) returned ok=false")
	}
	if eventMetric.EventType() != EventTypeCustom {
		t.Fatalf("event EventType() = %d, want %d", eventMetric.EventType(), EventTypeCustom)
	}
	if eventMetric.IsSessionDimension() {
		t.Fatal("event should not be a session dimension")
	}

	if _, ok := ParseMetricType("bad"); ok {
		t.Fatal("ParseMetricType(bad) returned ok=true")
	}
}

func TestNormalizeMetricLimit(t *testing.T) {
	tests := []struct {
		name  string
		limit int
		want  int
	}{
		{name: "zero", limit: 0, want: DefaultMetricLimit},
		{name: "negative", limit: -1, want: DefaultMetricLimit},
		{name: "too high", limit: MaxMetricLimit + 1, want: DefaultMetricLimit},
		{name: "valid", limit: 25, want: 25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizeMetricLimit(tt.limit); got != tt.want {
				t.Fatalf("NormalizeMetricLimit(%d) = %d, want %d", tt.limit, got, tt.want)
			}
		})
	}
}

func TestMetricTypeValues(t *testing.T) {
	want := []string{"path", "referrer", "browser", "os", "device", "country", "event"}
	if got := MetricTypeValues(); !reflect.DeepEqual(got, want) {
		t.Fatalf("MetricTypeValues() = %#v, want %#v", got, want)
	}
}
