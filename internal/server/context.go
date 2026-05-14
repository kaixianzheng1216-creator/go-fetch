package server

import (
	"context"
	"net/http"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
)

type contextKey string

const (
	userContextKey    contextKey = "user"
	requestContextKey contextKey = "request"
)

func (a *App) currentUser(ctx context.Context) (domain.User, bool) {
	userID := a.sessions.GetString(ctx, sessionUserIDKey)

	if userID == "" {
		return domain.User{}, false
	}

	user, err := a.store.GetUserByID(ctx, userID)

	if err != nil {
		return domain.User{}, false
	}

	return user, true
}

func withUser(ctx context.Context, user domain.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func userFromContext(ctx context.Context) domain.User {
	user, _ := ctx.Value(userContextKey).(domain.User)
	return user
}

func requestFromContext(ctx context.Context) *http.Request {
	r, _ := ctx.Value(requestContextKey).(*http.Request)
	return r
}

func captureRequest(ctx huma.Context, next func(huma.Context)) {
	r, _ := humachi.Unwrap(ctx)
	next(huma.WithContext(ctx, context.WithValue(ctx.Context(), requestContextKey, r)))
}
