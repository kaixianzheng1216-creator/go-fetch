package domain

import (
	"testing"
	"time"
)

func TestDateRangeUsesProvidedValues(t *testing.T) {
	startAt := int64(0)
	endAt := int64(3600000)

	start, end, unit := DateRange(&startAt, &endAt, string(DateUnitDay))

	if !start.Equal(time.UnixMilli(0)) {
		t.Fatalf("start = %s", start)
	}
	if !end.Equal(time.UnixMilli(3600000)) {
		t.Fatalf("end = %s", end)
	}
	if unit != DateUnitDay {
		t.Fatalf("unit = %q", unit)
	}
}

func TestMetricTypeEventType(t *testing.T) {
	if MetricTypeEvent.EventType() != EventTypeCustom {
		t.Fatalf("event metric should use custom event type")
	}
	if MetricTypePath.EventType() != EventTypePageView {
		t.Fatalf("path metric should use pageview event type")
	}
}
