package httpapi

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/cors"
	"github.com/google/uuid"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/session"
)

type requestContextKey struct{}

type userContextKey struct{}

func (srv server) collectCORSMiddleware(next http.Handler) http.Handler {
	corsHandler := cors.Handler(cors.Options{
		AllowedOrigins: srv.config.CollectCORSAllowedOrigins,
		AllowedMethods: []string{http.MethodPost, http.MethodOptions},
		AllowedHeaders: []string{"Content-Type"},
		MaxAge:         300,
	})(next)

	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/api/collect" {
			corsHandler.ServeHTTP(responseWriter, request)
			return
		}

		next.ServeHTTP(responseWriter, request)
	})
}

func captureRequest(ctx huma.Context, next func(huma.Context)) {
	request, _ := humachi.Unwrap(ctx)
	if request == nil {
		next(ctx)
		return
	}

	next(huma.WithContext(ctx, withRequest(ctx.Context(), request)))
}

func (srv server) requireAuth(humaAPI huma.API) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		user, isAuthenticated, err := srv.currentSessionUser(ctx.Context())
		if err != nil {
			if err := huma.WriteErr(humaAPI, ctx, http.StatusInternalServerError, errorMessageCurrentUserLoadFailed); err != nil {
				slog.Debug("write current user error", "error", err)
			}
			return
		}
		if !isAuthenticated {
			if err := huma.WriteErr(humaAPI, ctx, http.StatusUnauthorized, errorMessageUnauthenticated); err != nil {
				slog.Debug("write unauthenticated error", "error", err)
			}
			return
		}

		next(huma.WithContext(ctx, withUser(ctx.Context(), user)))
	}
}

func (srv server) currentSessionUser(ctx context.Context) (domain.User, bool, error) {
	if srv.sessions == nil {
		return domain.User{}, false, nil
	}

	userIDValue := srv.sessions.GetString(ctx, session.UserIDKey)
	if userIDValue == "" {
		return domain.User{}, false, nil
	}

	userID, err := uuid.Parse(userIDValue)
	if err != nil {
		return domain.User{}, false, nil
	}

	user, err := srv.users.GetByID(ctx, userID)
	if err != nil {
		if isNotFound(err) {
			return domain.User{}, false, nil
		}
		return domain.User{}, false, err
	}

	return user, true, nil
}

func currentUser(ctx context.Context) (domain.User, bool) {
	user, _ := ctx.Value(userContextKey{}).(domain.User)
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
	request, _ := ctx.Value(requestContextKey{}).(*http.Request)
	return request
}

func withRequest(ctx context.Context, request *http.Request) context.Context {
	return context.WithValue(ctx, requestContextKey{}, request)
}

func withUser(ctx context.Context, user domain.User) context.Context {
	return context.WithValue(ctx, userContextKey{}, user)
}
