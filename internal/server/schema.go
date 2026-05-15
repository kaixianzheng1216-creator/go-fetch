package server

import (
	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
)

func (CollectionType) Schema(huma.Registry) *huma.Schema {
	return &huma.Schema{
		Type: huma.TypeString,
		Enum: enumValues(domain.CollectionTypeValues()),
	}
}

type dateUnitParam string

func (dateUnitParam) Schema(huma.Registry) *huma.Schema {
	return &huma.Schema{
		Type:    huma.TypeString,
		Enum:    enumValues(domain.DateUnitValues()),
		Default: string(domain.DefaultDateUnit),
	}
}

type metricTypeParam string

func (metricTypeParam) Schema(huma.Registry) *huma.Schema {
	return &huma.Schema{
		Type: huma.TypeString,
		Enum: enumValues(domain.MetricTypeValues()),
	}
}

type metricLimit int

func (metricLimit) Schema(huma.Registry) *huma.Schema {
	minValue := 1.0
	maxValue := float64(domain.MaxMetricLimit)
	return &huma.Schema{
		Type:    huma.TypeInteger,
		Default: domain.DefaultMetricLimit,
		Minimum: &minValue,
		Maximum: &maxValue,
	}
}

func optionalTimeParam(value int64) *int64 {
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
