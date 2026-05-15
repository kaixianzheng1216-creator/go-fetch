package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/handler"
)

func RegisterAuth(api huma.API, authHandler handler.AuthHandler, authMiddleware huma.Middlewares) {
	huma.Register(
		api,
		NewOperation(
			http.MethodPost,
			"/api/login",
			"login",
			"登录",
			"Auth",
			http.StatusBadRequest,
			http.StatusUnauthorized,
			http.StatusUnprocessableEntity,
			http.StatusInternalServerError,
		),
		authHandler.Login,
	)

	huma.Register(
		api,
		NewOperation(
			http.MethodPost,
			"/api/logout",
			"logout",
			"退出登录",
			"Auth",
			http.StatusInternalServerError,
		),
		authHandler.Logout,
	)

	huma.Register(
		api,
		WithAuth(
			NewOperation(
				http.MethodGet,
				"/api/me",
				"getCurrentUser",
				"获取当前用户",
				"Auth",
				http.StatusUnauthorized,
			),
			authMiddleware,
		),
		authHandler.CurrentUser,
	)
}
