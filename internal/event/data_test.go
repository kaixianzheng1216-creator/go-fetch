package event

import (
	"math"
	"reflect"
	"sort"
	"testing"
	"time"
)

func TestFlattenEventData(t *testing.T) {
	createdAt := time.Date(2026, 5, 15, 9, 30, 0, 0, time.UTC)
	tests := []struct {
		name string
		data map[string]any
		want []FlatEventData
	}{
		{
			name: "flattens supported values",
			data: map[string]any{
				"createdAt": "2026-05-15T09:30:00Z",
				"meta": map[string]any{
					"ok":    true,
					"score": float64(12.5),
				},
				"tags": []any{"alpha", "beta"},
			},
			want: []FlatEventData{
				{
					Key:         "createdAt",
					StringValue: "2026-05-15T09:30:00Z",
					DateValue:   ptr(createdAt),
					DataType:    EventDataTypeDate,
				},
				{
					Key:         "meta.ok",
					StringValue: "true",
					DataType:    EventDataTypeBoolean,
				},
				{
					Key:         "meta.score",
					StringValue: "12.5",
					NumberValue: ptr(12.5),
					DataType:    EventDataTypeNumber,
				},
				{
					Key:         "tags",
					StringValue: `["alpha","beta"]`,
					DataType:    EventDataTypeArray,
				},
			},
		},
		{
			name: "skips non-finite numbers",
			data: map[string]any{
				"inf": math.Inf(1),
				"nan": math.NaN(),
			},
			want: nil,
		},
		{
			name: "stores nil as empty string value",
			data: map[string]any{"value": nil},
			want: []FlatEventData{{Key: "value", DataType: EventDataTypeString}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FlattenEventData(tt.data)
			sortFlatEventData(got)
			sortFlatEventData(tt.want)

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("FlattenEventData() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestTruncateEventDataValue(t *testing.T) {
	tests := []struct {
		name  string
		value string
		max   int
		want  string
	}{
		{name: "keeps shorter value", value: "abc", max: 5, want: "abc"},
		{name: "truncates by rune count", value: "abcdef", max: 3, want: "abc"},
		{name: "empty for non-positive max", value: "abcdef", max: 0, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := truncateEventDataValue(tt.value, tt.max); got != tt.want {
				t.Fatalf("truncateEventDataValue(%q, %d) = %q, want %q", tt.value, tt.max, got, tt.want)
			}
		})
	}
}

func sortFlatEventData(items []FlatEventData) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Key < items[j].Key
	})
}

func ptr[T any](value T) *T {
	return &value
}
