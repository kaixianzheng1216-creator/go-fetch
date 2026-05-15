package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/handler"
)

func registerAuthRoutes(api huma.API, authHandler handler.AuthHandler, authMiddleware huma.Middlewares) {
	huma.Register(
		api,
		huma.Operation{
			Method:      http.MethodPost,
			Path:        "/api/login",
			OperationID: "login",
			Summary:     "登录",
			Tags:        []string{"Auth"},
		},
		authHandler.Login,
	)

	huma.Register(
		api,
		huma.Operation{
			Method:      http.MethodPost,
			Path:        "/api/logout",
			OperationID: "logout",
			Summary:     "退出登录",
			Tags:        []string{"Auth"},
		},
		authHandler.Logout,
	)

	huma.Register(
		api,
		huma.Operation{
			Method:      http.MethodGet,
			Path:        "/api/me",
			OperationID: "getCurrentUser",
			Summary:     "获取当前用户",
			Tags:        []string{"Auth"},
			Security:    []map[string][]string{{"sessionCookie": {}}},
			Middlewares: authMiddleware,
		},
		authHandler.CurrentUser,
	)
}
