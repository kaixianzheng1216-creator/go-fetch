package domain

import (
	"errors"
	"time"
)

type CollectionType string

const (
	CollectionTypeEvent CollectionType = "event"
)

func ParseCollectionType(value string) (CollectionType, bool) {
	if value == "" {
		return CollectionTypeEvent, true
	}

	switch CollectionType(value) {
	case CollectionTypeEvent:
		return CollectionTypeEvent, true
	default:
		return "", false
	}
}

func CollectionTypeValues() []string {
	return []string{string(CollectionTypeEvent)}
}

type DateUnit string

const (
	DateUnitHour  DateUnit = "hour"
	DateUnitDay   DateUnit = "day"
	DateUnitMonth DateUnit = "month"

	DefaultDateUnit     = DateUnitHour
	DefaultDateLookback = 24 * time.Hour
)

func ParseDateUnit(value string) DateUnit {
	switch DateUnit(value) {
	case DateUnitMonth:
		return DateUnitMonth
	case DateUnitDay:
		return DateUnitDay
	default:
		return DefaultDateUnit
	}
}

func DateUnitValues() []string {
	return []string{string(DateUnitHour), string(DateUnitDay), string(DateUnitMonth)}
}

func DateRange(startAt, endAt *int64, unit string) (time.Time, time.Time, DateUnit) {
	now := time.Now()
	start := now.Add(-DefaultDateLookback)
	end := now

	if startAt != nil {
		start = time.UnixMilli(*startAt)
	}

	if endAt != nil {
		end = time.UnixMilli(*endAt)
	}

	return start, end, ParseDateUnit(unit)
}

func DateTruncUnit(unit DateUnit) string {
	return string(ParseDateUnit(string(unit)))
}

func FormatBucket(t time.Time, unit DateUnit) string {
	switch ParseDateUnit(string(unit)) {
	case DateUnitMonth:
		return t.Format("2006-01")
	case DateUnitDay:
		return t.Format("01-02")
	default:
		return t.Format("15:04")
	}
}

type MetricType string

const (
	MetricTypePath     MetricType = "path"
	MetricTypeReferrer MetricType = "referrer"
	MetricTypeBrowser  MetricType = "browser"
	MetricTypeOS       MetricType = "os"
	MetricTypeDevice   MetricType = "device"
	MetricTypeCountry  MetricType = "country"
	MetricTypeEvent    MetricType = "event"

	DefaultMetricLimit = 10
	MaxMetricLimit     = 100
)

var ErrUnsupportedMetricType = errors.New("unsupported metric type")

func ParseMetricType(value string) (MetricType, bool) {
	switch MetricType(value) {
	case MetricTypePath, MetricTypeReferrer, MetricTypeBrowser, MetricTypeOS, MetricTypeDevice, MetricTypeCountry, MetricTypeEvent:
		return MetricType(value), true
	default:
		return "", false
	}
}

func MetricTypeValues() []string {
	return []string{
		string(MetricTypePath),
		string(MetricTypeReferrer),
		string(MetricTypeBrowser),
		string(MetricTypeOS),
		string(MetricTypeDevice),
		string(MetricTypeCountry),
		string(MetricTypeEvent),
	}
}

func (m MetricType) EventType() EventType {
	if m == MetricTypeEvent {
		return EventTypeCustom
	}

	return EventTypePageView
}

func NormalizeMetricLimit(limit int) int {
	if limit <= 0 || limit > MaxMetricLimit {
		return DefaultMetricLimit
	}

	return limit
}
