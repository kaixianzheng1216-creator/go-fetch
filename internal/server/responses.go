package server

import (
	"encoding/json"
	"reflect"

	"go-fetch/internal/httpapi"

	"github.com/danielgtaylor/huma/v2"
)

type jsonBody[T any] struct {
	Body T
}

type collectResponseBody struct {
	value any
}

func collectResultBody(result httpapi.CollectResult) collectResponseBody {
	return collectResponseBody{value: result}
}

func collectOKBody() collectResponseBody {
	return collectResponseBody{value: httpapi.OKString{OK: "true"}}
}

func (b collectResponseBody) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.value)
}

func (b collectResponseBody) Schema(registry huma.Registry) *huma.Schema {
	return &huma.Schema{
		OneOf: []*huma.Schema{
			huma.SchemaFromType(registry, reflect.TypeOf(httpapi.CollectResult{})),
			huma.SchemaFromType(registry, reflect.TypeOf(httpapi.OKString{})),
		},
	}
}
