package server

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

func (a *App) requireHumaAuth(api huma.API) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		user, ok := a.currentUser(ctx.Context())
		if !ok {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "unauthorized")
			return
		}
		next(huma.WithContext(ctx, withUser(ctx.Context(), user)))
	}
}
