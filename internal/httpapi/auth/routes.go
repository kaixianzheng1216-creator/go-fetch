package auth

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/httpapi"
)

func Register(api huma.API, handler Handler, authMiddleware huma.Middlewares) {
	loginOp := httpapi.NewOperation(
		http.MethodPost,
		"/api/login",
		"login",
		"Auth",
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusUnprocessableEntity,
		http.StatusInternalServerError,
	)

	huma.Register(api, loginOp, handler.Login)

	logoutOp := httpapi.NewOperation(
		http.MethodPost,
		"/api/logout",
		"logout",
		"Auth",
		http.StatusInternalServerError,
	)

	huma.Register(api, logoutOp, handler.Logout)

	meOp := httpapi.NewOperation(
		http.MethodGet,
		"/api/me",
		"getCurrentUser",
		"Auth",
		http.StatusUnauthorized,
	)

	huma.Register(api, httpapi.WithAuth(meOp, authMiddleware), handler.CurrentUser)
}
