package server

import (
	"reflect"

	"go-fetch/internal/domain"
	"go-fetch/internal/httpapi"

	"github.com/danielgtaylor/huma/v2"
)

type optionalParam[T any] struct {
	Value T
	IsSet bool
}

func (v optionalParam[T]) Schema(registry huma.Registry) *huma.Schema {
	return huma.SchemaFromType(registry, reflect.TypeOf(v.Value))
}

func (v *optionalParam[T]) Receiver() reflect.Value {
	return reflect.ValueOf(v).Elem().FieldByName("Value")
}

func (v *optionalParam[T]) OnParamSet(isSet bool, _ any) {
	v.IsSet = isSet
}

func optionalValuePtr[T any](value optionalParam[T]) *T {
	if !value.IsSet {
		return nil
	}
	return &value.Value
}

type dateUnitParam string

func (dateUnitParam) Schema(huma.Registry) *huma.Schema {
	return httpapi.StringEnumSchema(domain.DateUnitValues(), string(domain.DefaultDateUnit))
}

type metricTypeParam string

func (metricTypeParam) Schema(huma.Registry) *huma.Schema {
	return httpapi.StringEnumSchema(domain.MetricTypeValues(), "")
}

type metricLimit int

func (metricLimit) Schema(huma.Registry) *huma.Schema {
	return httpapi.IntRangeSchema(domain.DefaultMetricLimit, 1, domain.MaxMetricLimit)
}
