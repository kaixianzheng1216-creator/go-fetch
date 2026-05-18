package domain

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
