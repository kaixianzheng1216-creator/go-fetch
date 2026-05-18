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

func DateTruncUnit(unit DateUnit) string {
	return string(ParseDateUnit(string(unit)))
}

func FormatBucket(bucketTime time.Time, unit DateUnit) string {
	switch ParseDateUnit(string(unit)) {
	case DateUnitMonth:
		return bucketTime.Format("2006-01")
	case DateUnitDay:
		return bucketTime.Format("01-02")
	default:
		return bucketTime.Format("15:04")
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

func (metricType MetricType) EventType() EventType {
	if metricType == MetricTypeEvent {
		return EventTypeCustom
	}

	return EventTypePageView
}

func (metricType MetricType) IsSessionDimension() bool {
	switch metricType {
	case MetricTypeBrowser, MetricTypeOS, MetricTypeDevice, MetricTypeCountry:
		return true
	default:
		return false
	}
}

func NormalizeMetricLimit(limit int) int {
	if limit <= 0 || limit > MaxMetricLimit {
		return DefaultMetricLimit
	}

	return limit
}
