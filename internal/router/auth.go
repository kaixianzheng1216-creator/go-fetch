package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/handler"
)

func registerAuthRoutes(api huma.API, authHandler handler.AuthHandler, authMiddleware huma.Middlewares) {
	huma.Register(
		api,
		operation(
			http.MethodPost,
			"/api/login",
			"login",
			"登录",
			"Auth",
		),
		authHandler.Login,
	)

	huma.Register(
		api,
		operation(
			http.MethodPost,
			"/api/logout",
			"logout",
			"退出登录",
			"Auth",
		),
		authHandler.Logout,
	)

	huma.Register(
		api,
		requireAuth(
			operation(
				http.MethodGet,
				"/api/me",
				"getCurrentUser",
				"获取当前用户",
				"Auth",
			),
			authMiddleware,
		),
		authHandler.CurrentUser,
	)
}
