package httpapi

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

func currentUser(ctx context.Context) domain.User {
	user, _ := ctx.Value(userContextKey).(domain.User)
	return user
}

func requestFromContext(ctx context.Context) *http.Request {
	request, _ := ctx.Value(requestContextKey).(*http.Request)
	return request
}

func withRequest(ctx context.Context, request *http.Request) context.Context {
	return context.WithValue(ctx, requestContextKey, request)
}

func withUser(ctx context.Context, user domain.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}
