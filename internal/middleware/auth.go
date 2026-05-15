package middleware

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
)

type CurrentUserFunc func(context.Context) (domain.User, bool, error)
type WithUserFunc func(context.Context, domain.User) context.Context
type WithRequestFunc func(context.Context, *http.Request) context.Context

func CaptureRequest(assignRequest WithRequestFunc) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		request, _ := humachi.Unwrap(ctx)
		if request == nil {
			next(ctx)
			return
		}

		next(huma.WithContext(ctx, assignRequest(ctx.Context(), request)))
	}
}

func RequireAuth(api huma.API, currentUser CurrentUserFunc, assignUser WithUserFunc) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		user, isAuthenticated, err := currentUser(ctx.Context())
		if err != nil {
			_ = huma.WriteErr(api, ctx, http.StatusInternalServerError, "加载当前用户失败")
			return
		}
		if !isAuthenticated {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "未登录")
			return
		}

		next(huma.WithContext(ctx, assignUser(ctx.Context(), user)))
	}
}
