package httpapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/service"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/session"
)

type loginInput struct {
	Body struct {
		Username string `json:"username" required:"true" minLength:"1"`
		Password string `json:"password" required:"true" minLength:"1" writeOnly:"true"`
	}
}

type loginOutput struct {
	Body struct {
		User struct {
			ID       uuid.UUID `json:"id" format:"uuid"`
			Username string    `json:"username"`
		} `json:"user"`
	}
}

type userOutput struct {
	Body struct {
		ID       uuid.UUID `json:"id" format:"uuid"`
		Username string    `json:"username"`
	}
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

	output := &loginOutput{}
	output.Body.User.ID = user.ID
	output.Body.User.Username = user.Username
	return output, nil
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

	output := &userOutput{}
	output.Body.ID = user.ID
	output.Body.Username = user.Username
	return output, nil
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

func loginError(err error) error {
	if errors.Is(err, service.ErrInvalidCredentials) {
		return huma.Error401Unauthorized(errorMessageInvalidCredentials)
	}
	return huma.Error500InternalServerError(errorMessageUserLoadFailed)
}
