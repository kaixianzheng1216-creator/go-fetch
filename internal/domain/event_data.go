package domain

import (
	"encoding/json"
	"fmt"
	"maps"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/textutil"
)

const (
	maxEventDataKeyLength   = 500
	maxEventDataValueLength = 500
)

type EventDataType int

const (
	EventDataTypeString  EventDataType = 1
	EventDataTypeNumber  EventDataType = 2
	EventDataTypeBoolean EventDataType = 3
	EventDataTypeDate    EventDataType = 4
	EventDataTypeArray   EventDataType = 5
)

type FlatEventData struct {
	Key         string
	StringValue string
	NumberValue *float64
	DateValue   *time.Time
	DataType    EventDataType
}

type eventDataFlattener struct {
	items []FlatEventData
}

func FlattenEventData(data map[string]any) []FlatEventData {
	flattener := eventDataFlattener{}
	for _, key := range eventDataKeys(data) {
		flattener.walk(key, data[key])
	}

	return flattener.items
}

func (flattener *eventDataFlattener) walk(prefix string, value any) {
	switch typedValue := value.(type) {
	case map[string]any:
		for _, key := range eventDataKeys(typedValue) {
			flattener.walk(joinEventDataKey(prefix, key), typedValue[key])
		}
	case []any:
		flattener.appendArray(prefix, typedValue)
	case bool:
		flattener.appendBool(prefix, typedValue)
	case string:
		flattener.appendString(prefix, typedValue)
	case nil:
		flattener.appendNull(prefix)
	default:
		if numberValue, ok := eventDataNumber(typedValue); ok {
			flattener.appendNumber(prefix, numberValue)
			return
		}
		flattener.appendText(prefix, fmt.Sprint(typedValue))
	}
}

func (flattener *eventDataFlattener) appendArray(key string, value []any) {
	bytes, err := json.Marshal(value)
	if err != nil {
		bytes = []byte(fmt.Sprint(value))
	}

	flattener.items = append(flattener.items, FlatEventData{
		Key:         eventDataKey(key),
		StringValue: textutil.TruncateRunes(string(bytes), maxEventDataValueLength),
		DataType:    EventDataTypeArray,
	})
}

func (flattener *eventDataFlattener) appendBool(key string, value bool) {
	flattener.items = append(flattener.items, FlatEventData{
		Key:         eventDataKey(key),
		StringValue: strconv.FormatBool(value),
		DataType:    EventDataTypeBoolean,
	})
}

func (flattener *eventDataFlattener) appendString(key, value string) {
	if dateValue, ok := parseEventDataTime(value); ok {
		flattener.appendDate(key, dateValue)
		return
	}

	flattener.appendText(key, value)
}

func (flattener *eventDataFlattener) appendText(key, value string) {
	flattener.items = append(flattener.items, FlatEventData{
		Key:         eventDataKey(key),
		StringValue: textutil.TruncateRunes(value, maxEventDataValueLength),
		DataType:    EventDataTypeString,
	})
}

func (flattener *eventDataFlattener) appendDate(key string, value time.Time) {
	value = value.UTC()
	flattener.items = append(flattener.items, FlatEventData{
		Key:         eventDataKey(key),
		StringValue: value.Format(time.RFC3339Nano),
		DateValue:   &value,
		DataType:    EventDataTypeDate,
	})
}

func (flattener *eventDataFlattener) appendNumber(key string, value float64) {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return
	}

	numberValue := value
	flattener.items = append(flattener.items, FlatEventData{
		Key:         eventDataKey(key),
		StringValue: strconv.FormatFloat(value, 'g', -1, 64),
		NumberValue: &numberValue,
		DataType:    EventDataTypeNumber,
	})
}

func (flattener *eventDataFlattener) appendNull(key string) {
	flattener.items = append(flattener.items, FlatEventData{
		Key:      eventDataKey(key),
		DataType: EventDataTypeString,
	})
}

func eventDataNumber(value any) (float64, bool) {
	switch typedValue := value.(type) {
	case float64:
		return typedValue, true
	case float32:
		return float64(typedValue), true
	case int:
		return float64(typedValue), true
	case int8:
		return float64(typedValue), true
	case int16:
		return float64(typedValue), true
	case int32:
		return float64(typedValue), true
	case int64:
		return float64(typedValue), true
	case uint:
		return float64(typedValue), true
	case uint8:
		return float64(typedValue), true
	case uint16:
		return float64(typedValue), true
	case uint32:
		return float64(typedValue), true
	case uint64:
		return float64(typedValue), true
	case json.Number:
		numberValue, err := typedValue.Float64()
		return numberValue, err == nil
	default:
		return 0, false
	}
}

func parseEventDataTime(value string) (time.Time, bool) {
	if !strings.Contains(value, "T") {
		return time.Time{}, false
	}

	for _, layout := range [...]string{
		time.RFC3339Nano,
		"2006-01-02T15:04:05.000",
		"2006-01-02T15:04:05",
	} {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			return parsed.UTC(), true
		}
	}

	return time.Time{}, false
}

func joinEventDataKey(prefix, key string) string {
	if prefix == "" {
		return key
	}

	return prefix + "." + key
}

func eventDataKey(key string) string {
	return textutil.TruncateRunes(key, maxEventDataKeyLength)
}

func eventDataKeys(data map[string]any) []string {
	return slices.Sorted(maps.Keys(data))
}
