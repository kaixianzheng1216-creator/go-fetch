package httpapi

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
)

type contextKey string

const (
	userContextKey    contextKey = "user"
	requestContextKey contextKey = "request"
)

func currentUser(ctx context.Context) (domain.User, bool) {
	user, _ := ctx.Value(userContextKey).(domain.User)
	return user, user.ID != uuid.Nil
}

func requireCurrentUser(ctx context.Context) (domain.User, error) {
	user, ok := currentUser(ctx)
	if !ok {
		return domain.User{}, huma.Error401Unauthorized(errorMessageUnauthenticated)
	}

	return user, nil
}

func currentUserID(ctx context.Context) (uuid.UUID, error) {
	user, err := requireCurrentUser(ctx)
	if err != nil {
		return uuid.Nil, err
	}

	return user.ID, nil
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
