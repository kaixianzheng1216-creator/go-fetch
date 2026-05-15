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
	tests := []struct {
		name   string
		metric MetricType
		want   EventType
	}{
		{name: "custom event metric", metric: MetricTypeEvent, want: EventTypeCustom},
		{name: "page path metric", metric: MetricTypePath, want: EventTypePageView},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.metric.EventType(); got != tt.want {
				t.Fatalf("EventType() = %d, want %d", got, tt.want)
			}
		})
	}
}
