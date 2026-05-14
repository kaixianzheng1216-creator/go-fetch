package domain

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
		switch v := value.(type) {
		case map[string]any:
			for key, child := range v {
				walk(joinEventDataKey(prefix, key), child)
			}
		case []any:
			bytes, err := json.Marshal(v)
			if err != nil {
				bytes = []byte(fmt.Sprint(v))
			}
			result = append(result, FlatEventData{
				Key:         prefix,
				StringValue: truncateEventDataValue(string(bytes), maxEventDataValueLength),
				DataType:    EventDataTypeArray,
			})
		case float64:
			if !math.IsNaN(v) && !math.IsInf(v, 0) {
				n := v
				result = append(result, FlatEventData{
					Key:         prefix,
					StringValue: fmt.Sprintf("%g", v),
					NumberValue: &n,
					DataType:    EventDataTypeNumber,
				})
			}
		case bool:
			result = append(result, FlatEventData{
				Key:         prefix,
				StringValue: strconv.FormatBool(v),
				DataType:    EventDataTypeBoolean,
			})
		case string:
			if dateValue, ok := parseEventDataTime(v); ok {
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
				StringValue: truncateEventDataValue(v, maxEventDataValueLength),
				DataType:    EventDataTypeString,
			})
		case nil:
			result = append(result, FlatEventData{Key: prefix, DataType: EventDataTypeString})
		default:
			result = append(result, FlatEventData{
				Key:         prefix,
				StringValue: truncateEventDataValue(fmt.Sprint(v), maxEventDataValueLength),
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
	for i := range value {
		if count == max {
			return value[:i]
		}

		count++
	}

	return value
}
