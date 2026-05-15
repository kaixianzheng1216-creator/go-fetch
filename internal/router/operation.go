package router

import "github.com/danielgtaylor/huma/v2"

func NewOperation(method, path, operationID, summary, tag string, statusCodes ...int) huma.Operation {
	return huma.Operation{
		Method:      method,
		Path:        path,
		OperationID: operationID,
		Summary:     summary,
		Tags:        []string{tag},
		Errors:      statusCodes,
	}
}

func WithDefaultStatus(op huma.Operation, status int) huma.Operation {
	op.DefaultStatus = status
	return op
}

func WithMaxBodyBytes(op huma.Operation, maxBodyBytes int64) huma.Operation {
	op.MaxBodyBytes = maxBodyBytes
	return op
}

func WithAuth(op huma.Operation, middlewares huma.Middlewares) huma.Operation {
	op.Security = []map[string][]string{{"sessionCookie": {}}}
	op.Middlewares = append(op.Middlewares, middlewares...)
	return op
}
