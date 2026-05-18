package httpapi

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/google/uuid"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/session"
)

func captureRequest(ctx huma.Context, next func(huma.Context)) {
	request, _ := humachi.Unwrap(ctx)
	if request == nil {
		next(ctx)
		return
	}

	next(huma.WithContext(ctx, withRequest(ctx.Context(), request)))
}

func (apiServer server) requireAuth(humaAPI huma.API) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		user, isAuthenticated, err := apiServer.currentSessionUser(ctx.Context())
		if err != nil {
			if err := huma.WriteErr(humaAPI, ctx, http.StatusInternalServerError, "加载当前用户失败"); err != nil {
				slog.Debug("write current user error", "error", err)
			}
			return
		}
		if !isAuthenticated {
			if err := huma.WriteErr(humaAPI, ctx, http.StatusUnauthorized, "未登录"); err != nil {
				slog.Debug("write unauthenticated error", "error", err)
			}
			return
		}

		next(huma.WithContext(ctx, withUser(ctx.Context(), user)))
	}
}

func (apiServer server) currentSessionUser(ctx context.Context) (domain.User, bool, error) {
	if apiServer.sessions == nil {
		return domain.User{}, false, nil
	}

	userIDValue := apiServer.sessions.GetString(ctx, session.UserIDKey)
	if userIDValue == "" {
		return domain.User{}, false, nil
	}

	userID, err := uuid.Parse(userIDValue)
	if err != nil {
		return domain.User{}, false, nil
	}

	user, err := apiServer.store.GetUserByID(ctx, userID)
	if err != nil {
		if isNotFound(err) {
			return domain.User{}, false, nil
		}
		return domain.User{}, false, err
	}

	return user, true, nil
}
