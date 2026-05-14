package domain

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

	if byKey["plan"].DataType != EventDataTypeString {
		t.Fatalf("plan DataType = %d", byKey["plan"].DataType)
	}
	if byKey["paid"].DataType != EventDataTypeBoolean {
		t.Fatalf("paid DataType = %d", byKey["paid"].DataType)
	}
	if byKey["amount"].DataType != EventDataTypeNumber || byKey["amount"].NumberValue == nil {
		t.Fatalf("amount = %#v", byKey["amount"])
	}
	if byKey["items"].DataType != EventDataTypeArray {
		t.Fatalf("items DataType = %d", byKey["items"].DataType)
	}
	if byKey["since"].DataType != EventDataTypeDate || byKey["since"].DateValue == nil {
		t.Fatalf("since = %#v", byKey["since"])
	}
}
