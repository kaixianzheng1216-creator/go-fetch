package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	requestdto "github.com/kaixianzheng1216-creator/go-fetch/internal/dto/request"
	responsedto "github.com/kaixianzheng1216-creator/go-fetch/internal/dto/response"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/service"
)

type Sessions interface {
	RenewToken(ctx context.Context) error
	Put(ctx context.Context, key string, val any)
	Destroy(ctx context.Context) error
}

type AuthHandler struct {
	auth        service.Auth
	sessions    Sessions
	userIDKey   string
	currentUser func(context.Context) domain.User
}

func NewAuth(
	auth service.Auth,
	sessions Sessions,
	userIDKey string,
	currentUser func(context.Context) domain.User,
) AuthHandler {
	return AuthHandler{auth: auth, sessions: sessions, userIDKey: userIDKey, currentUser: currentUser}
}

type loginRequest struct {
	Body requestdto.LoginRequest
}

func (handler AuthHandler) Login(ctx context.Context, input *loginRequest) (*responsedto.LoginOutput, error) {
	user, err := handler.auth.Login(ctx, input.Body.Username, input.Body.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return nil, huma.Error401Unauthorized("用户名或密码错误")
		}
		return nil, huma.Error500InternalServerError("加载用户失败")
	}

	if err := handler.startUserSession(ctx, user.ID); err != nil {
		return nil, huma.Error500InternalServerError("创建登录会话失败")
	}

	return responsedto.NewLoginOutput(responsedto.LoginResponse{User: responsedto.ToUser(user)}), nil
}

func (handler AuthHandler) Logout(ctx context.Context, _ *requestdto.Empty) (*responsedto.OKOutput, error) {
	if err := handler.sessions.Destroy(ctx); err != nil {
		return nil, huma.Error500InternalServerError("退出登录失败")
	}

	return responsedto.NewOKOutput(), nil
}

func (handler AuthHandler) CurrentUser(ctx context.Context, _ *requestdto.Empty) (*responsedto.UserOutput, error) {
	return responsedto.NewUserOutput(responsedto.ToUser(handler.currentUser(ctx))), nil
}

func (handler AuthHandler) startUserSession(ctx context.Context, userID string) error {
	if err := handler.sessions.RenewToken(ctx); err != nil {
		return fmt.Errorf("renew session token: %w", err)
	}

	handler.sessions.Put(ctx, handler.userIDKey, userID)
	return nil
}
