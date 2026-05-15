package summary

import (
	"github.com/danielgtaylor/huma/v2"

	eventdomain "github.com/kaixianzheng1216-creator/go-fetch/internal/domain/event"
)

type DateUnitParam string

func (DateUnitParam) Schema(huma.Registry) *huma.Schema {
	return &huma.Schema{
		Type:    huma.TypeString,
		Enum:    enumValues(eventdomain.DateUnitValues()),
		Default: string(eventdomain.DefaultDateUnit),
	}
}

type MetricTypeParam string

func (MetricTypeParam) Schema(huma.Registry) *huma.Schema {
	return &huma.Schema{
		Type: huma.TypeString,
		Enum: enumValues(eventdomain.MetricTypeValues()),
	}
}

type MetricLimit int

func (MetricLimit) Schema(huma.Registry) *huma.Schema {
	minValue := 1.0
	maxValue := float64(eventdomain.MaxMetricLimit)
	return &huma.Schema{
		Type:    huma.TypeInteger,
		Default: eventdomain.DefaultMetricLimit,
		Minimum: &minValue,
		Maximum: &maxValue,
	}
}

func OptionalTimeParam(value int64) *int64 {
	if value == 0 {
		return nil
	}
	return &value
}

func enumValues(values []string) []any {
	result := make([]any, 0, len(values))
	for _, value := range values {
		result = append(result, value)
	}
	return result
}
