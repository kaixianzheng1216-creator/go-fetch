package server

import (
	"context"
	"net/http"

	userdomain "github.com/kaixianzheng1216-creator/go-fetch/internal/domain/user"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/session"
)

type contextKey string

const (
	userContextKey    contextKey = "user"
	requestContextKey contextKey = "request"
)

func withUser(ctx context.Context, user userdomain.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func userFromContext(ctx context.Context) userdomain.User {
	user, _ := ctx.Value(userContextKey).(userdomain.User)

	return user
}

func withRequest(ctx context.Context, request *http.Request) context.Context {
	return context.WithValue(ctx, requestContextKey, request)
}

func requestFromContext(ctx context.Context) *http.Request {
	request, _ := ctx.Value(requestContextKey).(*http.Request)

	return request
}

func (a *App) currentUser(ctx context.Context) (userdomain.User, bool, error) {
	userID := a.sessions.GetString(ctx, session.UserIDKey)

	if userID == "" {
		return userdomain.User{}, false, nil
	}

	user, err := a.store.GetUserByID(ctx, userID)

	if err != nil {
		if isNotFound(err) {
			return userdomain.User{}, false, nil
		}

		return userdomain.User{}, false, err
	}

	return user, true, nil
}
