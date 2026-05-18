package domain

import "slices"

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

var metricTypes = [...]MetricType{
	MetricTypePath,
	MetricTypeReferrer,
	MetricTypeBrowser,
	MetricTypeOS,
	MetricTypeDevice,
	MetricTypeCountry,
	MetricTypeEvent,
}

type Metric struct {
	Name     string
	Views    int64
	Visitors int64
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
	}

	return false
}

func NormalizeMetricLimit(limit int) int {
	if limit <= 0 || limit > MaxMetricLimit {
		return DefaultMetricLimit
	}

	return limit
}
