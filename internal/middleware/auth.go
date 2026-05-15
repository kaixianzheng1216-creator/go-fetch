package middleware

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	userdomain "github.com/kaixianzheng1216-creator/go-fetch/internal/user"
)

type CurrentUserFunc func(context.Context) (userdomain.User, bool, error)
type WithUserFunc func(context.Context, userdomain.User) context.Context
type WithRequestFunc func(context.Context, *http.Request) context.Context

func CaptureRequest(withRequest WithRequestFunc) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		request, _ := humachi.Unwrap(ctx)

		if request == nil {
			next(ctx)
			return
		}

		next(huma.WithContext(ctx, withRequest(ctx.Context(), request)))
	}
}

func RequireAuth(api huma.API, currentUser CurrentUserFunc, withUser WithUserFunc) func(huma.Context, func(huma.Context)) {
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

		next(huma.WithContext(ctx, withUser(ctx.Context(), user)))
	}
}
