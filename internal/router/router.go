package router

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/handler"
	servermiddleware "github.com/kaixianzheng1216-creator/go-fetch/internal/middleware"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/repository"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/service"
	webassets "github.com/kaixianzheng1216-creator/go-fetch/web"
)

const (
	sessionCookieName = "go_fetch_session"
	userIDSessionKey  = "user_id"

	contentTypeHTML        = "text/html; charset=utf-8"
	contentTypeJS          = "application/javascript; charset=utf-8"
	contentTypeProblemJSON = "application/problem+json"
)

type contextKey string

const (
	userContextKey    contextKey = "user"
	requestContextKey contextKey = "request"
)

func New(dataStore *repository.Store) http.Handler {
	sessions := newSessionManager(dataStore)
	chiRouter := chi.NewRouter()

	chiRouter.Use(chimiddleware.RealIP)
	chiRouter.Use(chimiddleware.RequestID)
	chiRouter.Use(chimiddleware.Recoverer)
	chiRouter.Use(chimiddleware.Logger)
	chiRouter.Use(chimiddleware.Timeout(60 * time.Second))
	chiRouter.Use(sessions.LoadAndSave)

	api := humachi.New(chiRouter, humaConfig())
	api.UseMiddleware(servermiddleware.CaptureRequest(withRequest))
	authMiddleware := huma.Middlewares{servermiddleware.RequireAuth(api, currentUser(dataStore, sessions), withUser)}
	registerAPI(
		api,
		handler.NewAuth(service.NewAuth(dataStore, isNotFound), sessions, userIDSessionKey, userFromContext),
		handler.NewCollect(service.NewCollect(dataStore), requestFromContext, isNotFound),
		handler.NewWebsite(service.NewWebsite(dataStore), userFromContext, websiteLookupError),
		handler.NewStats(service.NewStats(dataStore), userFromContext, websiteLookupError),
		authMiddleware,
	)

	chiRouter.Get("/assets/*", http.FileServer(http.FS(webassets.DashboardFS())).ServeHTTP)
	chiRouter.Get("/script.js", handleScript)
	chiRouter.Get("/*", handleFrontend)

	return chiRouter
}

func OpenAPIJSON() ([]byte, error) {
	chiRouter := chi.NewRouter()
	api := humachi.New(chiRouter, humaConfig())

	registerAPI(api, handler.AuthHandler{}, handler.CollectHandler{}, handler.WebsiteHandler{}, handler.StatsHandler{}, nil)

	return json.MarshalIndent(api.OpenAPI(), "", "  ")
}

func registerAPI(
	api huma.API,
	authHandler handler.AuthHandler,
	collectHandler handler.CollectHandler,
	websiteHandler handler.WebsiteHandler,
	statsHandler handler.StatsHandler,
	authMiddleware huma.Middlewares,
) {
	registerCollectRoutes(api, collectHandler)
	registerAuthRoutes(api, authHandler, authMiddleware)
	registerWebsiteRoutes(api, websiteHandler, authMiddleware)
	registerStatsRoutes(api, statsHandler, authMiddleware)
}

func humaConfig() huma.Config {
	config := huma.DefaultConfig("go-fetch Analytics API", "0.1.0")
	config.DocsPath = "/api/docs"
	config.SchemasPath = ""
	config.CreateHooks = nil
	config.Servers = []*huma.Server{{URL: "/"}}
	config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"sessionCookie": {
			Type: "apiKey",
			In:   "cookie",
			Name: sessionCookieName,
		},
	}
	return config
}

func operation(method, path, operationID, summary, tag string) huma.Operation {
	return huma.Operation{
		Method:      method,
		Path:        path,
		OperationID: operationID,
		Summary:     summary,
		Tags:        []string{tag},
	}
}

func requireAuth(operation huma.Operation, middlewares huma.Middlewares) huma.Operation {
	operation.Security = []map[string][]string{{"sessionCookie": {}}}
	operation.Middlewares = append(operation.Middlewares, middlewares...)
	return operation
}

func newSessionManager(dataStore *repository.Store) *scs.SessionManager {
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

func currentUser(dataStore *repository.Store, sessions *scs.SessionManager) servermiddleware.CurrentUserFunc {
	return func(ctx context.Context) (domain.User, bool, error) {
		userID := sessions.GetString(ctx, userIDSessionKey)
		if userID == "" {
			return domain.User{}, false, nil
		}

		user, err := dataStore.GetUserByID(ctx, userID)
		if err != nil {
			if isNotFound(err) {
				return domain.User{}, false, nil
			}
			return domain.User{}, false, err
		}

		return user, true, nil
	}
}

func withRequest(ctx context.Context, request *http.Request) context.Context {
	return context.WithValue(ctx, requestContextKey, request)
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

func isNotFound(err error) bool {
	return errors.Is(err, repository.ErrNotFound)
}

func websiteLookupError(err error) error {
	if err == nil {
		return nil
	}
	if isNotFound(err) {
		return huma.Error404NotFound("站点不存在")
	}
	return huma.Error500InternalServerError("加载站点失败")
}

func handleScript(responseWriter http.ResponseWriter, _ *http.Request) {
	script, err := webassets.TrackerScript()
	if err != nil {
		http.Error(responseWriter, "tracking script is missing", http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Set("Content-Type", contentTypeJS)
	_, _ = responseWriter.Write(script)
}

func handleFrontend(responseWriter http.ResponseWriter, request *http.Request) {
	switch {
	case request.Method != http.MethodGet:
		http.NotFound(responseWriter, request)
		return
	case strings.HasPrefix(request.URL.Path, "/api/"):
		responseWriter.Header().Set("Content-Type", contentTypeProblemJSON)
		responseWriter.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(responseWriter).Encode(huma.ErrorModel{
			Title:  "接口不存在",
			Status: http.StatusNotFound,
			Detail: "接口不存在",
		})
		return
	}

	indexHTML, err := webassets.IndexHTML()
	if err != nil {
		http.Error(responseWriter, "dashboard build is missing", http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Set("Content-Type", contentTypeHTML)
	_, _ = responseWriter.Write(indexHTML)
}
