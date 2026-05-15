package event

import "testing"

func TestFlattenEventDataKeepsValueTypes(t *testing.T) {
	data := map[string]any{
		"plan":   "pro",
		"paid":   true,
		"amount": float64(12.5),
		"items":  []any{"a", "b"},
		"since":  "2026-05-14T12:30:00Z",
	}

	items := FlattenEventData(data)
	byKey := make(map[string]FlatEventData, len(items))
	for _, item := range items {
		byKey[item.Key] = item
	}

	tests := []struct {
		key      string
		dataType EventDataType
		hasNum   bool
		hasDate  bool
	}{
		{key: "plan", dataType: EventDataTypeString},
		{key: "paid", dataType: EventDataTypeBoolean},
		{key: "amount", dataType: EventDataTypeNumber, hasNum: true},
		{key: "items", dataType: EventDataTypeArray},
		{key: "since", dataType: EventDataTypeDate, hasDate: true},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			item, ok := byKey[tt.key]
			if !ok {
				t.Fatalf("missing key %q", tt.key)
			}
			if item.DataType != tt.dataType {
				t.Fatalf("DataType = %d, want %d", item.DataType, tt.dataType)
			}
			if (item.NumberValue != nil) != tt.hasNum {
				t.Fatalf("NumberValue present = %t, want %t", item.NumberValue != nil, tt.hasNum)
			}
			if (item.DateValue != nil) != tt.hasDate {
				t.Fatalf("DateValue present = %t, want %t", item.DateValue != nil, tt.hasDate)
			}
		})
	}
}
