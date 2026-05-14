package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"golang.org/x/crypto/bcrypt"
)

func registerAuthRoutes(api huma.API, app *App, auth huma.Middlewares) {
	loginOp := newOperation(
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

	logoutOp := newOperation(
		http.MethodPost,
		"/api/logout",
		"logout",
		"Auth",
		http.StatusInternalServerError,
	)

	huma.Register(api, logoutOp, app.logout)

	meOp := newOperation(
		http.MethodGet,
		"/api/me",
		"getCurrentUser",
		"Auth",
		http.StatusUnauthorized,
	)

	huma.Register(api, withAuth(meOp, auth), app.me)
}

func (a *App) login(ctx context.Context, input *loginInput) (*jsonBody[LoginResponse], error) {
	user, err := a.store.GetUserByUsername(ctx, input.Body.Username)
	if err != nil {
		if isNotFound(err) {
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

	response := LoginResponse{
		User: toUser(user),
	}

	return jsonResponse(response), nil
}

func (a *App) logout(ctx context.Context, _ *emptyInput) (*jsonBody[OK], error) {
	if err := a.sessions.Destroy(ctx); err != nil {
		return nil, huma.Error500InternalServerError("退出登录失败")
	}

	return okResponse(), nil
}

func (a *App) me(ctx context.Context, _ *emptyInput) (*jsonBody[User], error) {
	response := toUser(userFromContext(ctx))

	return jsonResponse(response), nil
}

func (a *App) startSession(ctx context.Context, userID string) error {
	if err := a.sessions.RenewToken(ctx); err != nil {
		return fmt.Errorf("刷新会话令牌失败: %w", err)
	}

	a.sessions.Put(ctx, sessionUserIDKey, userID)

	return nil
}
