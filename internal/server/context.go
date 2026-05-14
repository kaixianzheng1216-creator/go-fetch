package server

import (
	"context"
	"net/http"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
)

type contextKey string

const (
	userContextKey    contextKey = "user"
	requestContextKey contextKey = "request"
)

func (a *App) currentUser(ctx context.Context) (domain.User, bool, error) {
	userID := a.sessions.GetString(ctx, sessionUserIDKey)

	if userID == "" {
		return domain.User{}, false, nil
	}

	user, err := a.store.GetUserByID(ctx, userID)

	if err != nil {
		if isNotFound(err) {
			return domain.User{}, false, nil
		}

		return domain.User{}, false, err
	}

	return user, true, nil
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
