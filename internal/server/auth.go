package server

import (
	"context"
	"fmt"
	"net/http"

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

func (a *App) login(ctx context.Context, input *loginInput) (*jsonBody[LoginResponse], error) {
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

	return &jsonBody[LoginResponse]{Body: LoginResponse{User: UserFromDomain(user)}}, nil
}

func (a *App) logout(ctx context.Context, _ *emptyInput) (*jsonBody[OK], error) {
	if err := a.sessions.Destroy(ctx); err != nil {
		return nil, huma.Error500InternalServerError("退出登录失败")
	}

	return &jsonBody[OK]{Body: OK{OK: true}}, nil
}

func (a *App) me(ctx context.Context, _ *emptyInput) (*jsonBody[User], error) {
	return &jsonBody[User]{Body: UserFromDomain(userFromContext(ctx))}, nil
}

func (a *App) startSession(ctx context.Context, userID string) error {
	if err := a.sessions.RenewToken(ctx); err != nil {
		return fmt.Errorf("刷新会话令牌失败: %w", err)
	}

	a.sessions.Put(ctx, sessionUserIDKey, userID)

	return nil
}

func (a *App) requireHumaAuth(api huma.API) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		user, ok := a.currentUser(ctx.Context())
		if !ok {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "未登录或登录已失效")
			return
		}

		next(huma.WithContext(ctx, withUser(ctx.Context(), user)))
	}
}
