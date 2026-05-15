package event

import (
	"reflect"
	"testing"
	"time"
)

func TestParseCollectionType(testRunner *testing.T) {
	tests := []struct {
		name       string
		value      string
		expected   CollectionType
		expectedOK bool
	}{
		{name: "empty defaults to event", value: "", expected: CollectionTypeEvent, expectedOK: true},
		{name: "event", value: "event", expected: CollectionTypeEvent, expectedOK: true},
		{name: "unsupported", value: "pageview", expectedOK: false},
	}

	for _, testCase := range tests {
		testRunner.Run(testCase.name, func(testRunner *testing.T) {
			actual, actualOK := ParseCollectionType(testCase.value)
			if actual != testCase.expected || actualOK != testCase.expectedOK {
				testRunner.Fatalf("ParseCollectionType(%q) = (%q, %v), want (%q, %v)", testCase.value, actual, actualOK, testCase.expected, testCase.expectedOK)
			}
		})
	}
}

func TestParseDateUnit(testRunner *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected DateUnit
	}{
		{name: "month", value: "month", expected: DateUnitMonth},
		{name: "day", value: "day", expected: DateUnitDay},
		{name: "hour", value: "hour", expected: DefaultDateUnit},
		{name: "invalid", value: "invalid", expected: DefaultDateUnit},
		{name: "empty", value: "", expected: DefaultDateUnit},
	}

	for _, testCase := range tests {
		testRunner.Run(testCase.name, func(testRunner *testing.T) {
			if actual := ParseDateUnit(testCase.value); actual != testCase.expected {
				testRunner.Fatalf("ParseDateUnit(%q) = %q, want %q", testCase.value, actual, testCase.expected)
			}
		})
	}
}

func TestDateRange(testRunner *testing.T) {
	startMillis := time.Date(2026, 5, 15, 8, 0, 0, 0, time.UTC).UnixMilli()
	endMillis := time.Date(2026, 5, 15, 9, 0, 0, 0, time.UTC).UnixMilli()

	start, end, unit := DateRange(&startMillis, &endMillis, "day")

	if !start.Equal(time.UnixMilli(startMillis)) {
		testRunner.Fatalf("start = %s, want %s", start, time.UnixMilli(startMillis))
	}
	if !end.Equal(time.UnixMilli(endMillis)) {
		testRunner.Fatalf("end = %s, want %s", end, time.UnixMilli(endMillis))
	}
	if unit != DateUnitDay {
		testRunner.Fatalf("unit = %q, want %q", unit, DateUnitDay)
	}
}

func TestFormatBucket(testRunner *testing.T) {
	bucket := time.Date(2026, 5, 15, 9, 30, 0, 0, time.UTC)
	tests := []struct {
		unit     DateUnit
		expected string
	}{
		{unit: DateUnitMonth, expected: "2026-05"},
		{unit: DateUnitDay, expected: "05-15"},
		{unit: DateUnitHour, expected: "09:30"},
	}

	for _, testCase := range tests {
		testRunner.Run(string(testCase.unit), func(testRunner *testing.T) {
			if actual := FormatBucket(bucket, testCase.unit); actual != testCase.expected {
				testRunner.Fatalf("FormatBucket(%s, %q) = %q, want %q", bucket, testCase.unit, actual, testCase.expected)
			}
		})
	}
}

func TestMetricRules(testRunner *testing.T) {
	metric, parsed := ParseMetricType("browser")
	if !parsed {
		testRunner.Fatal("ParseMetricType(browser) returned ok=false")
	}

	if metric.EventType() != EventTypePageView {
		testRunner.Fatalf("browser EventType() = %d, want %d", metric.EventType(), EventTypePageView)
	}
	if !metric.IsSessionDimension() {
		testRunner.Fatal("browser should be a session dimension")
	}

	eventMetric, parsed := ParseMetricType("event")
	if !parsed {
		testRunner.Fatal("ParseMetricType(event) returned ok=false")
	}
	if eventMetric.EventType() != EventTypeCustom {
		testRunner.Fatalf("event EventType() = %d, want %d", eventMetric.EventType(), EventTypeCustom)
	}
	if eventMetric.IsSessionDimension() {
		testRunner.Fatal("event should not be a session dimension")
	}

	if _, parsed := ParseMetricType("bad"); parsed {
		testRunner.Fatal("ParseMetricType(bad) returned ok=true")
	}
}

func TestNormalizeMetricLimit(testRunner *testing.T) {
	tests := []struct {
		name     string
		limit    int
		expected int
	}{
		{name: "zero", limit: 0, expected: DefaultMetricLimit},
		{name: "negative", limit: -1, expected: DefaultMetricLimit},
		{name: "too high", limit: MaxMetricLimit + 1, expected: DefaultMetricLimit},
		{name: "valid", limit: 25, expected: 25},
	}

	for _, testCase := range tests {
		testRunner.Run(testCase.name, func(testRunner *testing.T) {
			if actual := NormalizeMetricLimit(testCase.limit); actual != testCase.expected {
				testRunner.Fatalf("NormalizeMetricLimit(%d) = %d, want %d", testCase.limit, actual, testCase.expected)
			}
		})
	}
}

func TestMetricTypeValues(testRunner *testing.T) {
	expected := []string{"path", "referrer", "browser", "os", "device", "country", "event"}
	if actual := MetricTypeValues(); !reflect.DeepEqual(actual, expected) {
		testRunner.Fatalf("MetricTypeValues() = %#v, want %#v", actual, expected)
	}
}
