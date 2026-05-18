package domain

import (
	"errors"
	"slices"
	"time"
)

type DateUnit string

const (
	DateUnitHour  DateUnit = "hour"
	DateUnitDay   DateUnit = "day"
	DateUnitMonth DateUnit = "month"

	DefaultDateUnit     = DateUnitHour
	DefaultDateLookback = 24 * time.Hour
)

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

var (
	ErrUnsupportedDateUnit   = errors.New("unsupported date unit")
	ErrUnsupportedMetricType = errors.New("unsupported metric type")
)

var (
	dateUnits = [...]DateUnit{
		DateUnitHour,
		DateUnitDay,
		DateUnitMonth,
	}

	metricTypes = [...]MetricType{
		MetricTypePath,
		MetricTypeReferrer,
		MetricTypeBrowser,
		MetricTypeOS,
		MetricTypeDevice,
		MetricTypeCountry,
		MetricTypeEvent,
	}
)

type WebsiteStats struct {
	Pageviews       int64
	Visitors        int64
	Visits          int64
	Bounces         int64
	TotalTime       int64
	AvgVisitSeconds int64
}

type PageviewBucket struct {
	Time     time.Time
	Label    string
	Views    int64
	Visitors int64
}

type Metric struct {
	Name     string
	Views    int64
	Visitors int64
}

func ParseDateUnit(value string) (DateUnit, bool) {
	if value == "" {
		return DefaultDateUnit, true
	}

	dateUnit := DateUnit(value)
	if slices.Contains(dateUnits[:], dateUnit) {
		return dateUnit, true
	}

	return "", false
}

func DateUnitValues() []string {
	return stringEnumValues(dateUnits[:])
}

func DateTruncUnit(unit DateUnit) string {
	return string(NormalizeDateUnit(unit))
}

func FormatBucket(bucketTime time.Time, unit DateUnit) string {
	switch NormalizeDateUnit(unit) {
	case DateUnitMonth:
		return bucketTime.Format("2006-01")
	case DateUnitDay:
		return bucketTime.Format("01-02")
	default:
		return bucketTime.Format("15:04")
	}
}

func NormalizeDateUnit(unit DateUnit) DateUnit {
	parsedUnit, ok := ParseDateUnit(string(unit))
	if !ok {
		return DefaultDateUnit
	}

	return parsedUnit
}

func ParseMetricType(value string) (MetricType, bool) {
	metricType := MetricType(value)
	if slices.Contains(metricTypes[:], metricType) {
		return metricType, true
	}

	return "", false
}

func MetricTypeValues() []string {
	return stringEnumValues(metricTypes[:])
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

func stringEnumValues[T ~string](values []T) []string {
	result := make([]string, len(values))
	for i, value := range values {
		result[i] = string(value)
	}
	return result
}
