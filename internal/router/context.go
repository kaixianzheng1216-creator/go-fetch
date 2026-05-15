package router

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

func withRequest(ctx context.Context, request *http.Request) context.Context {
	return context.WithValue(ctx, requestContextKey, request)
}

func (server *Server) currentUser(ctx context.Context) (domain.User, bool, error) {
	userID := server.sessions.GetString(ctx, userIDSessionKey)
	if userID == "" {
		return domain.User{}, false, nil
	}

	user, err := server.store.GetUserByID(ctx, userID)
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
	request, _ := ctx.Value(requestContextKey).(*http.Request)
	return request
}
