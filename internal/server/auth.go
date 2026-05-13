package server

import (
	"context"
	"net/http"
	"time"

	authpkg "github.com/kaixianzheng1216-creator/go-fetch/internal/auth"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/httpapi"

	"github.com/danielgtaylor/huma/v2"
	"github.com/go-chi/httprate"
)

func registerAuthRoutes(api huma.API, app *App, auth huma.Middlewares) {
	loginOp := operation(http.MethodPost, "/api/login", "login", "Auth", http.StatusBadRequest, http.StatusUnauthorized, http.StatusInternalServerError)
	loginOp.SkipValidateBody = true
	if app != nil {
		loginOp.Middlewares = append(loginOp.Middlewares, adaptHTTPMiddleware(httprate.LimitByRealIP(10, time.Minute)))
	}
	huma.Register(api, loginOp, app.login)

	huma.Register(api, operation(http.MethodPost, "/api/logout", "logout", "Auth", http.StatusInternalServerError), app.logout)

	huma.Register(api, authenticated(operation(http.MethodGet, "/api/me", "getCurrentUser", "Auth", http.StatusUnauthorized), auth), app.me)
}

func (a *App) login(ctx context.Context, input *loginInput) (*jsonBody[httpapi.LoginResponse], error) {
	if input.Body.Username == "" || input.Body.Password == "" {
		return nil, huma.Error400BadRequest("username and password are required")
	}

	user, err := a.store.GetUserByUsername(ctx, input.Body.Username)
	if err != nil {
		if isStoreNotFound(err) {
			return nil, huma.Error401Unauthorized("incorrect username or password")
		}

		return nil, huma.Error500InternalServerError("failed to load user")
	}

	if !authpkg.CheckPassword(user.PasswordHash, input.Body.Password) {
		return nil, huma.Error401Unauthorized("incorrect username or password")
	}

	if err := a.startSession(ctx, user.ID); err != nil {
		return nil, huma.Error500InternalServerError("failed to create session")
	}

	return &jsonBody[httpapi.LoginResponse]{Body: httpapi.LoginResponse{User: httpapi.UserFromDomain(user)}}, nil
}

func (a *App) logout(ctx context.Context, _ *emptyInput) (*jsonBody[httpapi.OK], error) {
	if err := a.sessions.Destroy(ctx); err != nil {
		return nil, huma.Error500InternalServerError("failed to destroy session")
	}

	return &jsonBody[httpapi.OK]{Body: httpapi.OK{OK: true}}, nil
}

func (a *App) me(ctx context.Context, _ *emptyInput) (*jsonBody[httpapi.User], error) {
	return &jsonBody[httpapi.User]{Body: httpapi.UserFromDomain(userFromContext(ctx))}, nil
}

func (a *App) startSession(ctx context.Context, userID string) error {
	if err := a.sessions.RenewToken(ctx); err != nil {
		return err
	}

	a.sessions.Put(ctx, sessionUserIDKey, userID)

	return nil
}
