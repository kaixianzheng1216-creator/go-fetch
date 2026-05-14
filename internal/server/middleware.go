package server

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const httpRequestTimeout = 30 * time.Second

func (a *App) useHTTPMiddleware(r chi.Router) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(httpRequestTimeout))
	r.Use(a.sessions.LoadAndSave)
}

func (a *App) useAPIMiddleware(api huma.API) huma.Middlewares {
	api.UseMiddleware(captureRequest)

	return huma.Middlewares{a.requireAuth(api)}
}

func withAuth(op huma.Operation, middlewares huma.Middlewares) huma.Operation {
	op.Security = []map[string][]string{{"sessionCookie": {}}}

	op.Middlewares = append(op.Middlewares, middlewares...)

	return op
}

func (a *App) requireAuth(api huma.API) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		user, ok, err := a.currentUser(ctx.Context())

		if err != nil {
			_ = huma.WriteErr(api, ctx, http.StatusInternalServerError, "加载当前用户失败")

			return
		}

		if !ok {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "未登录或登录已失效")

			return
		}

		next(huma.WithContext(ctx, withUser(ctx.Context(), user)))
	}
}

func captureRequest(ctx huma.Context, next func(huma.Context)) {
	request, _ := humachi.Unwrap(ctx)

	if request == nil {
		next(ctx)

		return
	}

	next(huma.WithContext(ctx, context.WithValue(ctx.Context(), requestContextKey, request)))
}
