package httpapi

import "github.com/danielgtaylor/huma/v2"

func NewOperation(method, path, operationID, tag string, statusCodes ...int) huma.Operation {
	return huma.Operation{
		Method:      method,
		Path:        path,
		OperationID: operationID,
		Tags:        []string{tag},
		Errors:      statusCodes,
	}
}

func WithAuth(op huma.Operation, middlewares huma.Middlewares) huma.Operation {
	op.Security = []map[string][]string{{"sessionCookie": {}}}
	op.Middlewares = append(op.Middlewares, middlewares...)
	return op
}
