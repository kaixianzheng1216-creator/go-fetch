package router

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/session"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/handler"
	servermiddleware "github.com/kaixianzheng1216-creator/go-fetch/internal/middleware"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/repository"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/service"
	webassets "github.com/kaixianzheng1216-creator/go-fetch/web"
)

func New(dataStore *repository.Store) http.Handler {
	sessions := session.NewManager(dataStore)
	chiRouter := chi.NewRouter()

	chiRouter.Use(chimiddleware.RealIP)
	chiRouter.Use(chimiddleware.RequestID)
	chiRouter.Use(chimiddleware.Recoverer)
	chiRouter.Use(chimiddleware.Logger)
	chiRouter.Use(chimiddleware.Timeout(60 * time.Second))
	chiRouter.Use(collectCORSMiddleware)
	chiRouter.Use(sessions.LoadAndSave)

	api := humachi.New(chiRouter, humaConfig())
	api.UseMiddleware(servermiddleware.CaptureRequest(withRequest))
	authMiddleware := huma.Middlewares{servermiddleware.RequireAuth(
		api,
		func(ctx context.Context) (domain.User, bool, error) {
			userID := sessions.GetString(ctx, session.UserIDKey)
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
		},
		withUser,
	)}

	registerAPI(
		api,
		handler.NewAuth(service.NewAuth(dataStore, isNotFound), sessions, session.UserIDKey, userFromContext),
		handler.NewCollect(service.NewCollect(dataStore), requestFromContext, isNotFound),
		handler.NewWebsite(service.NewWebsite(dataStore), userFromContext, websiteLookupError),
		handler.NewStats(service.NewStats(dataStore), userFromContext, websiteLookupError),
		authMiddleware,
	)

	chiRouter.Get("/assets/*", http.FileServer(http.FS(webassets.DashboardFS())).ServeHTTP)
	chiRouter.Get("/script.js", func(responseWriter http.ResponseWriter, _ *http.Request) {
		script, err := webassets.TrackerScript()
		if err != nil {
			http.Error(responseWriter, "tracking script is missing", http.StatusInternalServerError)
			return
		}

		responseWriter.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		_, _ = responseWriter.Write(script)
	})
	chiRouter.Get("/*", func(responseWriter http.ResponseWriter, request *http.Request) {
		switch {
		case request.Method != http.MethodGet:
			http.NotFound(responseWriter, request)
			return
		case strings.HasPrefix(request.URL.Path, "/api/"):
			responseWriter.Header().Set("Content-Type", "application/problem+json")
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

		responseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = responseWriter.Write(indexHTML)
	})

	return chiRouter
}

type contextKey string

const (
	userContextKey    contextKey = "user"
	requestContextKey contextKey = "request"
)

func collectCORSMiddleware(next http.Handler) http.Handler {
	corsHandler := cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
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

func isNotFound(err error) bool {
	return errors.Is(err, repository.ErrNotFound)
}

func userFromContext(ctx context.Context) domain.User {
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

func websiteLookupError(err error) error {
	if err == nil {
		return nil
	}
	if isNotFound(err) {
		return huma.Error404NotFound("站点不存在")
	}
	return huma.Error500InternalServerError("加载站点失败")
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
			Name: session.CookieName,
		},
	}
	return config
}
