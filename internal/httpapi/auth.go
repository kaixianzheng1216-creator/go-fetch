package httpapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/service"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/session"
)

type loginInput struct {
	Body LoginRequest
}

type LoginRequest struct {
	Username string `json:"username" required:"true" minLength:"1"`
	Password string `json:"password" required:"true" minLength:"1" writeOnly:"true"`
}

type UserResponse struct {
	ID        uuid.UUID  `json:"id" format:"uuid"`
	Username  string     `json:"username"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

type LoginResponse struct {
	User UserResponse `json:"user"`
}

type loginOutput struct {
	Body LoginResponse
}

type userOutput struct {
	Body UserResponse
}

func (apiServer server) registerAuthRoutes(humaAPI huma.API, authMiddleware huma.Middlewares) {
	huma.Register(
		humaAPI,
		huma.Operation{
			Method:      http.MethodPost,
			Path:        "/api/login",
			OperationID: "login",
			Summary:     "登录",
			Tags:        []string{"Auth"},
		},
		apiServer.login,
	)

	huma.Register(
		humaAPI,
		huma.Operation{
			Method:      http.MethodPost,
			Path:        "/api/logout",
			OperationID: "logout",
			Summary:     "退出登录",
			Tags:        []string{"Auth"},
		},
		apiServer.logout,
	)

	huma.Register(
		humaAPI,
		huma.Operation{
			Method:      http.MethodGet,
			Path:        "/api/me",
			OperationID: "getCurrentUser",
			Summary:     "获取当前用户",
			Tags:        []string{"Auth"},
			Security:    []map[string][]string{{"sessionCookie": {}}},
			Middlewares: authMiddleware,
		},
		apiServer.getCurrentUser,
	)
}

func (apiServer server) login(ctx context.Context, input *loginInput) (*loginOutput, error) {
	user, err := apiServer.auth.Login(ctx, input.Body.Username, input.Body.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return nil, huma.Error401Unauthorized("用户名或密码错误")
		}
		return nil, huma.Error500InternalServerError("加载用户失败")
	}

	if err := apiServer.startUserSession(ctx, user.ID); err != nil {
		return nil, huma.Error500InternalServerError("创建登录会话失败")
	}

	return &loginOutput{Body: LoginResponse{User: toUserResponse(user)}}, nil
}

func (apiServer server) logout(ctx context.Context, _ *emptyInput) (*okOutput, error) {
	if apiServer.sessions == nil {
		return toOKOutput(), nil
	}
	if err := apiServer.sessions.Destroy(ctx); err != nil {
		return nil, huma.Error500InternalServerError("退出登录失败")
	}

	return toOKOutput(), nil
}

func (apiServer server) getCurrentUser(ctx context.Context, _ *emptyInput) (*userOutput, error) {
	return &userOutput{Body: toUserResponse(currentUser(ctx))}, nil
}

func (apiServer server) startUserSession(ctx context.Context, userID uuid.UUID) error {
	if apiServer.sessions == nil {
		return fmt.Errorf("session manager is not configured")
	}
	if err := apiServer.sessions.RenewToken(ctx); err != nil {
		return fmt.Errorf("renew session token: %w", err)
	}

	apiServer.sessions.Put(ctx, session.UserIDKey, userID.String())
	return nil
}

func toUserResponse(user domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}
}
