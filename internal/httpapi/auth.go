package httpapi

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
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

func (srv server) registerAuthRoutes(humaAPI huma.API, authMiddleware huma.Middlewares) {
	huma.Register(
		humaAPI,
		publicOperation(http.MethodPost, "/api/login", "login", "Log in", "Auth"),
		srv.login,
	)

	huma.Register(
		humaAPI,
		publicOperation(http.MethodPost, "/api/logout", "logout", "Log out", "Auth"),
		srv.logout,
	)

	huma.Register(
		humaAPI,
		securedOperation(http.MethodGet, "/api/me", "getCurrentUser", "Get current user", "Auth", authMiddleware),
		srv.getCurrentUser,
	)
}

func (srv server) login(ctx context.Context, input *loginInput) (*loginOutput, error) {
	user, err := srv.auth.Login(ctx, input.Body.Username, input.Body.Password)
	if err != nil {
		return nil, loginError(err)
	}

	if err := srv.startUserSession(ctx, user.ID); err != nil {
		return nil, huma.Error500InternalServerError(errorMessageLoginSessionCreate)
	}

	return &loginOutput{Body: LoginResponse{User: newUserResponse(user)}}, nil
}

func (srv server) logout(ctx context.Context, _ *emptyInput) (*okOutput, error) {
	if srv.sessions == nil {
		return newOKOutput(), nil
	}
	if err := srv.sessions.Destroy(ctx); err != nil {
		return nil, huma.Error500InternalServerError(errorMessageLogoutFailed)
	}

	return newOKOutput(), nil
}

func (srv server) getCurrentUser(ctx context.Context, _ *emptyInput) (*userOutput, error) {
	user, err := requireCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	return &userOutput{Body: newUserResponse(user)}, nil
}

func (srv server) startUserSession(ctx context.Context, userID uuid.UUID) error {
	if srv.sessions == nil {
		return fmt.Errorf("session manager is not configured")
	}
	if err := srv.sessions.RenewToken(ctx); err != nil {
		return fmt.Errorf("renew session token: %w", err)
	}

	srv.sessions.Put(ctx, session.UserIDKey, userID.String())
	return nil
}

func newUserResponse(user domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}
}
