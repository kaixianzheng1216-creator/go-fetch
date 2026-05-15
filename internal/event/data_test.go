package event

import (
	"math"
	"reflect"
	"sort"
	"testing"
	"time"
)

func TestFlattenEventData(testRunner *testing.T) {
	createdAt := time.Date(2026, 5, 15, 9, 30, 0, 0, time.UTC)
	tests := []struct {
		name     string
		data     map[string]any
		expected []FlatEventData
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
			expected: []FlatEventData{
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
			expected: nil,
		},
		{
			name:     "stores nil as empty string value",
			data:     map[string]any{"value": nil},
			expected: []FlatEventData{{Key: "value", DataType: EventDataTypeString}},
		},
	}

	for _, testCase := range tests {
		testRunner.Run(testCase.name, func(testRunner *testing.T) {
			actual := FlattenEventData(testCase.data)
			sortFlatEventData(actual)
			sortFlatEventData(testCase.expected)

			if !reflect.DeepEqual(actual, testCase.expected) {
				testRunner.Fatalf("FlattenEventData() = %#v, want %#v", actual, testCase.expected)
			}
		})
	}
}

func TestTruncateEventDataValue(testRunner *testing.T) {
	tests := []struct {
		name     string
		value    string
		max      int
		expected string
	}{
		{name: "keeps shorter value", value: "abc", max: 5, expected: "abc"},
		{name: "truncates by rune count", value: "abcdef", max: 3, expected: "abc"},
		{name: "empty for non-positive max", value: "abcdef", max: 0, expected: ""},
	}

	for _, testCase := range tests {
		testRunner.Run(testCase.name, func(testRunner *testing.T) {
			if actual := truncateEventDataValue(testCase.value, testCase.max); actual != testCase.expected {
				testRunner.Fatalf("truncateEventDataValue(%q, %d) = %q, want %q", testCase.value, testCase.max, actual, testCase.expected)
			}
		})
	}
}

func sortFlatEventData(items []FlatEventData) {
	sort.Slice(items, func(leftIndex, rightIndex int) bool {
		return items[leftIndex].Key < items[rightIndex].Key
	})
}

func ptr[T any](value T) *T {
	return &value
}
