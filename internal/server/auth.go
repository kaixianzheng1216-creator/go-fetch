package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/httpapi"

	"github.com/danielgtaylor/huma/v2"
	"golang.org/x/crypto/bcrypt"
)

func registerAuthRoutes(api huma.API, app *App, auth huma.Middlewares) {
	loginOp := operation(
		http.MethodPost,
		"/api/login",
		"login",
		"Auth",
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusUnprocessableEntity,
		http.StatusInternalServerError,
	)

	huma.Register(api, loginOp, app.login)

	logoutOp := operation(
		http.MethodPost,
		"/api/logout",
		"logout",
		"Auth",
		http.StatusInternalServerError,
	)

	huma.Register(api, logoutOp, app.logout)

	meOp := operation(
		http.MethodGet,
		"/api/me",
		"getCurrentUser",
		"Auth",
		http.StatusUnauthorized,
	)

	huma.Register(api, authenticated(meOp, auth), app.me)
}

func (a *App) login(ctx context.Context, input *loginInput) (*jsonBody[httpapi.LoginResponse], error) {
	user, err := a.store.GetUserByUsername(ctx, input.Body.Username)
	if err != nil {
		if isStoreNotFound(err) {
			return nil, huma.Error401Unauthorized("用户名或密码不正确")
		}

		return nil, huma.Error500InternalServerError("加载用户失败")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Body.Password)) != nil {
		return nil, huma.Error401Unauthorized("用户名或密码不正确")
	}

	if err := a.startSession(ctx, user.ID); err != nil {
		return nil, huma.Error500InternalServerError("创建登录会话失败")
	}

	return &jsonBody[httpapi.LoginResponse]{Body: httpapi.LoginResponse{User: httpapi.UserFromDomain(user)}}, nil
}

func (a *App) logout(ctx context.Context, _ *emptyInput) (*jsonBody[httpapi.OK], error) {
	if err := a.sessions.Destroy(ctx); err != nil {
		return nil, huma.Error500InternalServerError("退出登录失败")
	}

	return &jsonBody[httpapi.OK]{Body: httpapi.OK{OK: true}}, nil
}

func (a *App) me(ctx context.Context, _ *emptyInput) (*jsonBody[httpapi.User], error) {
	return &jsonBody[httpapi.User]{Body: httpapi.UserFromDomain(userFromContext(ctx))}, nil
}

func (a *App) startSession(ctx context.Context, userID string) error {
	if err := a.sessions.RenewToken(ctx); err != nil {
		return fmt.Errorf("刷新会话令牌失败: %w", err)
	}

	a.sessions.Put(ctx, sessionUserIDKey, userID)

	return nil
}
