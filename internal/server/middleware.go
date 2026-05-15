package server

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"

	userdomain "github.com/kaixianzheng1216-creator/go-fetch/internal/user"
)

type currentUserFunc func(context.Context) (userdomain.User, bool, error)
type withUserFunc func(context.Context, userdomain.User) context.Context
type withRequestFunc func(context.Context, *http.Request) context.Context

func captureRequestMiddleware(assignRequest withRequestFunc) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		request, _ := humachi.Unwrap(ctx)

		if request == nil {
			next(ctx)
			return
		}

		next(huma.WithContext(ctx, assignRequest(ctx.Context(), request)))
	}
}

func requireAuthMiddleware(api huma.API, currentUser currentUserFunc, assignUser withUserFunc) func(huma.Context, func(huma.Context)) {
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
