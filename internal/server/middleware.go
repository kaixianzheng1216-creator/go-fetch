package server

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/unrolled/secure"
)

func (a *App) secureHeaders() *secure.Secure {
	return secure.New(secure.Options{
		BrowserXssFilter:   true,
		ContentTypeNosniff: true,
		FrameDeny:          true,
		ReferrerPolicy:     "strict-origin-when-cross-origin",
		PermissionsPolicy:  "camera=(), microphone=(), geolocation=()",
		SSLProxyHeaders:    map[string]string{"X-Forwarded-Proto": "https"},
		STSSeconds:         stsSeconds(a.cfg.CookieSecure),
	})
}

func adaptHTTPMiddleware(mw func(http.Handler) http.Handler) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		r, w := humachi.Unwrap(ctx)
		mw(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
			next(huma.WithContext(ctx, r.Context()))
		})).ServeHTTP(w, r)
	}
}

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

func stsSeconds(enabled bool) int64 {
	if !enabled {
		return 0
	}
	return 31536000
}
