package model

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

const maxEventDataValueLength = 500

type FlatEventData struct {
	Key         string
	StringValue string
	NumberValue *float64
	DateValue   *time.Time
	DataType    EventDataType
}

func FlattenEventData(data map[string]any) []FlatEventData {
	var result []FlatEventData

	var walk func(prefix string, value any)
	walk = func(prefix string, value any) {
		switch typedValue := value.(type) {
		case map[string]any:
			for key, child := range typedValue {
				walk(joinEventDataKey(prefix, key), child)
			}
		case []any:
			bytes, err := json.Marshal(typedValue)
			if err != nil {
				bytes = []byte(fmt.Sprint(typedValue))
			}
			result = append(result, FlatEventData{
				Key:         prefix,
				StringValue: truncateEventDataValue(string(bytes), maxEventDataValueLength),
				DataType:    EventDataTypeArray,
			})
		case float64:
			if !math.IsNaN(typedValue) && !math.IsInf(typedValue, 0) {
				numberValue := typedValue
				result = append(result, FlatEventData{
					Key:         prefix,
					StringValue: fmt.Sprintf("%g", typedValue),
					NumberValue: &numberValue,
					DataType:    EventDataTypeNumber,
				})
			}
		case bool:
			result = append(result, FlatEventData{
				Key:         prefix,
				StringValue: strconv.FormatBool(typedValue),
				DataType:    EventDataTypeBoolean,
			})
		case string:
			if dateValue, hasDateValue := parseEventDataTime(typedValue); hasDateValue {
				result = append(result, FlatEventData{
					Key:         prefix,
					StringValue: dateValue.UTC().Format(time.RFC3339Nano),
					DateValue:   &dateValue,
					DataType:    EventDataTypeDate,
				})
				break
			}

			result = append(result, FlatEventData{
				Key:         prefix,
				StringValue: truncateEventDataValue(typedValue, maxEventDataValueLength),
				DataType:    EventDataTypeString,
			})
		case nil:
			result = append(result, FlatEventData{Key: prefix, DataType: EventDataTypeString})
		default:
			result = append(result, FlatEventData{
				Key:         prefix,
				StringValue: truncateEventDataValue(fmt.Sprint(typedValue), maxEventDataValueLength),
				DataType:    EventDataTypeString,
			})
		}
	}

	for key, value := range data {
		walk(key, value)
	}

	return result
}

func parseEventDataTime(value string) (time.Time, bool) {
	if !strings.Contains(value, "T") {
		return time.Time{}, false
	}

	for _, layout := range []string{time.RFC3339Nano, "2006-01-02T15:04:05.000", "2006-01-02T15:04:05"} {
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

func truncateEventDataValue(value string, max int) string {
	if max <= 0 {
		return ""
	}

	count := 0
	for index := range value {
		if count == max {
			return value[:index]
		}

		count++
	}

	return value
}
