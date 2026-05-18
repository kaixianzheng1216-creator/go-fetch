package repository

import (
	"encoding/json"
	"fmt"
	"maps"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/util"
)

const (
	maxEventDataKeyLength   = 500
	maxEventDataValueLength = 500
)

type eventDataType int

const (
	eventDataTypeString  eventDataType = 1
	eventDataTypeNumber  eventDataType = 2
	eventDataTypeBoolean eventDataType = 3
	eventDataTypeDate    eventDataType = 4
	eventDataTypeArray   eventDataType = 5
)

type eventDataRow struct {
	Key         string
	StringValue string
	NumberValue *float64
	DateValue   *time.Time
	DataType    eventDataType
}

func flattenEventData(data map[string]any) []eventDataRow {
	rows := make([]eventDataRow, 0, len(data))
	for _, key := range eventDataKeys(data) {
		rows = flattenEventDataValue(rows, key, data[key])
	}

	return rows
}

func flattenEventDataValue(rows []eventDataRow, key string, value any) []eventDataRow {
	if nested, ok := value.(map[string]any); ok {
		for _, childKey := range eventDataKeys(nested) {
			rows = flattenEventDataValue(rows, joinEventDataKey(key, childKey), nested[childKey])
		}
		return rows
	}

	row, ok := newEventDataRow(key, value)
	if !ok {
		return rows
	}
	return append(rows, row)
}

func newEventDataRow(key string, value any) (eventDataRow, bool) {
	key = eventDataKey(key)

	switch typedValue := value.(type) {
	case []any:
		return eventDataArrayRow(key, typedValue), true
	case bool:
		return eventDataTextRow(key, strconv.FormatBool(typedValue), eventDataTypeBoolean), true
	case string:
		if dateValue, ok := parseEventDataTime(typedValue); ok {
			return eventDataDateRow(key, dateValue), true
		}
		return eventDataTextRow(key, typedValue, eventDataTypeString), true
	case nil:
		return eventDataRow{Key: key, DataType: eventDataTypeString}, true
	default:
		numberValue, ok := eventDataNumber(typedValue)
		if !ok {
			return eventDataTextRow(key, fmt.Sprint(typedValue), eventDataTypeString), true
		}
		return eventDataNumberRow(key, numberValue)
	}
}

func eventDataTextRow(key, value string, dataType eventDataType) eventDataRow {
	return eventDataRow{
		Key:         key,
		StringValue: util.TruncateRunes(value, maxEventDataValueLength),
		DataType:    dataType,
	}
}

func eventDataArrayRow(key string, value []any) eventDataRow {
	bytes, err := json.Marshal(value)
	if err != nil {
		return eventDataTextRow(key, fmt.Sprint(value), eventDataTypeArray)
	}

	return eventDataTextRow(key, string(bytes), eventDataTypeArray)
}

func eventDataNumberRow(key string, value float64) (eventDataRow, bool) {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return eventDataRow{}, false
	}

	numberValue := value
	return eventDataRow{
		Key:         key,
		StringValue: strconv.FormatFloat(value, 'g', -1, 64),
		NumberValue: &numberValue,
		DataType:    eventDataTypeNumber,
	}, true
}

func eventDataDateRow(key string, value time.Time) eventDataRow {
	value = value.UTC()
	return eventDataRow{
		Key:         key,
		StringValue: value.Format(time.RFC3339Nano),
		DateValue:   &value,
		DataType:    eventDataTypeDate,
	}
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
	return util.TruncateRunes(key, maxEventDataKeyLength)
}

func eventDataKeys(data map[string]any) []string {
	return slices.Sorted(maps.Keys(data))
}
