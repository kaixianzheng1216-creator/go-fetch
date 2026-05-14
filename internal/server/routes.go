package server

import (
	"github.com/danielgtaylor/huma/v2"
)

func registerAPIRoutes(api huma.API, app *App) {
	api.UseMiddleware(captureRequest)

	auth := huma.Middlewares{app.requireHumaAuth(api)}
	registerCollectRoutes(api, app)
	registerAuthRoutes(api, app, auth)
	registerWebsiteRoutes(api, app, auth)
	registerAnalyticsRoutes(api, app, auth)
}

func operation(method, path, operationID, tag string, errors ...int) huma.Operation {
	return huma.Operation{
		Method:      method,
		Path:        path,
		OperationID: operationID,
		Tags:        []string{tag},
		Errors:      errors,
	}
}

func authenticated(op huma.Operation, middlewares huma.Middlewares) huma.Operation {
	op.Security = []map[string][]string{{"sessionCookie": {}}}

	op.Middlewares = append(op.Middlewares, middlewares...)

	return op
}
