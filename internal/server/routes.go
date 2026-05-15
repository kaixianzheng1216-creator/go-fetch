package server

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/httpapi/analytics"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/httpapi/auth"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/httpapi/events"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/httpapi/websites"
	assets "github.com/kaixianzheng1216-creator/go-fetch/internal/static"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/store"
	userdomain "github.com/kaixianzheng1216-creator/go-fetch/internal/user"
)

const (
	contentTypeHTML = "text/html; charset=utf-8"
	contentTypeJS   = "application/javascript; charset=utf-8"

	apiPrefix = "/api/"

	sessionCookieName = "go_fetch_session"
	userIDSessionKey  = "user_id"
)

type currentUserFunc func(context.Context) (userdomain.User, bool, error)
type withUserFunc func(context.Context, userdomain.User) context.Context
type withRequestFunc func(context.Context, *http.Request) context.Context

func newSessionManager(dataStore *store.Store) *scs.SessionManager {
	sessionManager := scs.New()

	sessionManager.Store = pgxstore.NewWithConfig(dataStore.Pool(), pgxstore.Config{
		TableName:       "app_sessions",
		CleanUpInterval: 10 * time.Minute,
	})

	sessionManager.Cookie.Name = sessionCookieName
	sessionManager.Cookie.Secure = true
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.SameSite = http.SameSiteLaxMode

	sessionManager.Lifetime = 24 * time.Hour

	return sessionManager
}

func (app *App) Routes() http.Handler {
	router := chi.NewRouter()

	router.Use(chimiddleware.RealIP)
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Timeout(60 * time.Second))
	router.Use(app.sessions.LoadAndSave)

	api := humachi.New(router, humaConfig())

	registerAPIRoutes(api, app)

	router.Get("/assets/*", app.handleFrontendAsset)
	router.Get("/script.js", app.handleScript)
	router.Get("/*", app.handleFrontend)

	return router
}

func registerAPIRoutes(api huma.API, app *App) {
	api.UseMiddleware(captureRequestMiddleware(withRequest))

	authMiddleware := huma.Middlewares{requireAuthMiddleware(api, app.currentUser, withUser)}

	authHandler := auth.New(app.store, app.sessions, userIDSessionKey, userFromContext, isNotFound)
	eventsHandler := events.New(app.store, requestFromContext, isNotFound)
	websitesHandler := websites.New(app.store, userFromContext, websiteLookupError)
	analyticsHandler := analytics.New(app.store, userFromContext, websiteLookupError)

	events.Register(api, eventsHandler)
	auth.Register(api, authHandler, authMiddleware)
	websites.Register(api, websitesHandler, authMiddleware)
	analytics.Register(api, analyticsHandler, authMiddleware)
}

func captureRequestMiddleware(assignRequest withRequestFunc) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		request, _ := humachi.Unwrap(ctx)

		if request == nil {
			next(ctx)
			return
		}

		next(huma.WithContext(ctx, assignRequest(ctx.Context(), request)))
	}
}

func requireAuthMiddleware(api huma.API, currentUser currentUserFunc, assignUser withUserFunc) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		user, isAuthenticated, err := currentUser(ctx.Context())

		if err != nil {
			_ = huma.WriteErr(api, ctx, http.StatusInternalServerError, "加载当前用户失败")
			return
		}

		if !isAuthenticated {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "未登录")
			return
		}

		next(huma.WithContext(ctx, assignUser(ctx.Context(), user)))
	}
}

func (app *App) handleFrontendAsset(responseWriter http.ResponseWriter, request *http.Request) {
	http.FileServer(http.FS(assets.DistFS())).ServeHTTP(responseWriter, request)
}

func (app *App) handleScript(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type", contentTypeJS)

	http.ServeFileFS(responseWriter, request, assets.StaticFS(), "script.js")
}

func (app *App) handleFrontend(responseWriter http.ResponseWriter, request *http.Request) {
	switch {
	case request.Method != http.MethodGet:
		http.NotFound(responseWriter, request)
		return

	case strings.HasPrefix(request.URL.Path, apiPrefix):
		writeProblemError(responseWriter, http.StatusNotFound, "接口不存在")
		return
	}

	indexHTML, err := assets.IndexHTML()
	if err != nil {
		http.Error(responseWriter, "前端构建产物不存在", http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Set("Content-Type", contentTypeHTML)

	_, _ = responseWriter.Write(indexHTML)
}
