package httpapi

import "github.com/danielgtaylor/huma/v2"

func StringEnumSchema(values []string, defaultValue string) *huma.Schema {
	schema := &huma.Schema{
		Type: huma.TypeString,
		Enum: enumValues(values),
	}
	if defaultValue != "" {
		schema.Default = defaultValue
	}
	return schema
}

func IntRangeSchema(defaultValue, minimum, maximum int) *huma.Schema {
	minValue := float64(minimum)
	maxValue := float64(maximum)
	return &huma.Schema{
		Type:    huma.TypeInteger,
		Default: defaultValue,
		Minimum: &minValue,
		Maximum: &maxValue,
	}
}

func enumValues(values []string) []any {
	result := make([]any, 0, len(values))
	for _, value := range values {
		result = append(result, value)
	}
	return result
}
