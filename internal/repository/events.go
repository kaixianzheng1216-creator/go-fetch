package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"math"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	storesqlc "github.com/kaixianzheng1216-creator/go-fetch/internal/repository/sqlc"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/textutil"
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

func (store *Store) SaveEvent(ctx context.Context, event domain.EventRecord) error {
	err := pgx.BeginFunc(ctx, store.pool, func(tx pgx.Tx) error {
		queries := store.queries.WithTx(tx)
		if err := queries.InsertSession(ctx, insertSessionParams(event)); err != nil {
			return fmt.Errorf("insert session: %w", err)
		}

		eventID := uuid.New()
		if err := queries.InsertEvent(ctx, insertEventParams(eventID, event)); err != nil {
			return fmt.Errorf("insert event: %w", err)
		}

		for _, item := range flattenEventData(event.Data) {
			if err := queries.InsertEventData(ctx, insertEventDataParams(event, eventID, item)); err != nil {
				return fmt.Errorf("insert event data: %w", err)
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("save event transaction: %w", err)
	}

	return nil
}

func insertSessionParams(event domain.EventRecord) storesqlc.InsertSessionParams {
	return storesqlc.InsertSessionParams{
		ID:         event.SessionID,
		WebsiteID:  event.WebsiteID,
		Browser:    event.Browser,
		OS:         event.OS,
		Device:     event.Device,
		Screen:     event.Screen,
		Language:   event.Language,
		Country:    event.Country,
		Region:     event.Region,
		City:       event.City,
		DistinctID: event.DistinctID,
		CreatedAt:  event.CreatedAt,
	}
}

func insertEventParams(eventID uuid.UUID, event domain.EventRecord) storesqlc.InsertEventParams {
	return storesqlc.InsertEventParams{
		ID:             eventID,
		WebsiteID:      event.WebsiteID,
		SessionID:      event.SessionID,
		VisitID:        event.VisitID,
		EventType:      int32(event.EventType),
		EventName:      event.EventName,
		URLPath:        event.URLPath,
		URLQuery:       event.URLQuery,
		ReferrerPath:   event.ReferrerPath,
		ReferrerQuery:  event.ReferrerQuery,
		ReferrerDomain: event.ReferrerDomain,
		PageTitle:      event.PageTitle,
		Hostname:       event.Hostname,
		UTMSource:      event.UTMSource,
		UTMMedium:      event.UTMMedium,
		UTMCampaign:    event.UTMCampaign,
		UTMContent:     event.UTMContent,
		UTMTerm:        event.UTMTerm,
		CreatedAt:      event.CreatedAt,
	}
}

func insertEventDataParams(event domain.EventRecord, eventID uuid.UUID, item eventDataRow) storesqlc.InsertEventDataParams {
	return storesqlc.InsertEventDataParams{
		ID:          uuid.New(),
		WebsiteID:   event.WebsiteID,
		EventID:     eventID,
		DataKey:     item.Key,
		StringValue: item.StringValue,
		NumberValue: item.NumberValue,
		DateValue:   item.DateValue,
		DataType:    int32(item.DataType),
		CreatedAt:   event.CreatedAt,
	}
}

func flattenEventData(data map[string]any) []eventDataRow {
	rows := make([]eventDataRow, 0, len(data))
	for _, key := range eventDataKeys(data) {
		rows = appendEventDataRows(rows, key, data[key])
	}

	return rows
}

func appendEventDataRows(rows []eventDataRow, prefix string, value any) []eventDataRow {
	switch typedValue := value.(type) {
	case map[string]any:
		for _, key := range eventDataKeys(typedValue) {
			rows = appendEventDataRows(rows, joinEventDataKey(prefix, key), typedValue[key])
		}
		return rows
	case []any:
		return appendEventDataArray(rows, prefix, typedValue)
	case bool:
		return appendEventDataBool(rows, prefix, typedValue)
	case string:
		return appendEventDataString(rows, prefix, typedValue)
	case nil:
		return appendEventDataNull(rows, prefix)
	default:
		if numberValue, ok := eventDataNumber(typedValue); ok {
			return appendEventDataNumber(rows, prefix, numberValue)
		}
		return appendEventDataText(rows, prefix, fmt.Sprint(typedValue))
	}
}

func appendEventDataArray(rows []eventDataRow, key string, value []any) []eventDataRow {
	bytes, err := json.Marshal(value)
	if err != nil {
		bytes = []byte(fmt.Sprint(value))
	}

	return append(rows, eventDataRow{
		Key:         eventDataKey(key),
		StringValue: textutil.TruncateRunes(string(bytes), maxEventDataValueLength),
		DataType:    eventDataTypeArray,
	})
}

func appendEventDataBool(rows []eventDataRow, key string, value bool) []eventDataRow {
	return append(rows, eventDataRow{
		Key:         eventDataKey(key),
		StringValue: strconv.FormatBool(value),
		DataType:    eventDataTypeBoolean,
	})
}

func appendEventDataString(rows []eventDataRow, key, value string) []eventDataRow {
	if dateValue, ok := parseEventDataTime(value); ok {
		return appendEventDataDate(rows, key, dateValue)
	}

	return appendEventDataText(rows, key, value)
}

func appendEventDataText(rows []eventDataRow, key, value string) []eventDataRow {
	return append(rows, eventDataRow{
		Key:         eventDataKey(key),
		StringValue: textutil.TruncateRunes(value, maxEventDataValueLength),
		DataType:    eventDataTypeString,
	})
}

func appendEventDataDate(rows []eventDataRow, key string, value time.Time) []eventDataRow {
	value = value.UTC()
	return append(rows, eventDataRow{
		Key:         eventDataKey(key),
		StringValue: value.Format(time.RFC3339Nano),
		DateValue:   &value,
		DataType:    eventDataTypeDate,
	})
}

func appendEventDataNumber(rows []eventDataRow, key string, value float64) []eventDataRow {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return rows
	}

	numberValue := value
	return append(rows, eventDataRow{
		Key:         eventDataKey(key),
		StringValue: strconv.FormatFloat(value, 'g', -1, 64),
		NumberValue: &numberValue,
		DataType:    eventDataTypeNumber,
	})
}

func appendEventDataNull(rows []eventDataRow, key string) []eventDataRow {
	return append(rows, eventDataRow{
		Key:      eventDataKey(key),
		DataType: eventDataTypeString,
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
