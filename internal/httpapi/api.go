package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/repository"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/service"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/session"
	webassets "github.com/kaixianzheng1216-creator/go-fetch/web"
)

type Config struct {
	CollectCORSAllowedOrigins []string
}

type sessionStore interface {
	RenewToken(ctx context.Context) error
	Put(ctx context.Context, key string, val any)
	Destroy(ctx context.Context) error
	GetString(ctx context.Context, key string) string
}

type server struct {
	store    *repository.Store
	sessions sessionStore
	auth     service.Auth
	collect  service.Collect
	websites service.Website
	stats    service.Stats
	config   Config
}

func New(store *repository.Store, sessions *scs.SessionManager, config Config) http.Handler {
	apiServer := server{
		store:    store,
		sessions: sessions,
		auth:     service.NewAuth(store, isNotFound),
		collect:  service.NewCollect(store),
		websites: service.NewWebsite(store),
		stats:    service.NewStats(store),
		config:   config.withDefaults(),
	}

	chiRouter := chi.NewRouter()
	chiRouter.Use(chimiddleware.RealIP)
	chiRouter.Use(chimiddleware.RequestID)
	chiRouter.Use(chimiddleware.Recoverer)
	chiRouter.Use(chimiddleware.Logger)
	chiRouter.Use(chimiddleware.Timeout(60 * time.Second))
	chiRouter.Use(apiServer.collectCORSMiddleware)
	if sessions != nil {
		chiRouter.Use(sessions.LoadAndSave)
	}

	humaAPI := humachi.New(chiRouter, humaConfig())
	humaAPI.UseMiddleware(captureRequest)
	apiServer.registerRoutes(humaAPI)
	apiServer.registerAssets(chiRouter)

	return chiRouter
}

func (config Config) withDefaults() Config {
	if len(config.CollectCORSAllowedOrigins) == 0 {
		config.CollectCORSAllowedOrigins = []string{"*"}
	}
	return config
}

func (apiServer server) registerRoutes(humaAPI huma.API) {
	authMiddleware := huma.Middlewares{apiServer.requireAuth(humaAPI)}

	apiServer.registerCollectRoutes(humaAPI)
	apiServer.registerAuthRoutes(humaAPI, authMiddleware)
	apiServer.registerWebsiteRoutes(humaAPI, authMiddleware)
	apiServer.registerStatsRoutes(humaAPI, authMiddleware)
}

func (apiServer server) registerAssets(chiRouter chi.Router) {
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
	chiRouter.Get("/*", spaHandler)
}

func (apiServer server) collectCORSMiddleware(next http.Handler) http.Handler {
	corsHandler := cors.Handler(cors.Options{
		AllowedOrigins: apiServer.config.CollectCORSAllowedOrigins,
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

func spaHandler(responseWriter http.ResponseWriter, request *http.Request) {
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

func enumValues(values []string) []any {
	result := make([]any, 0, len(values))
	for _, value := range values {
		result = append(result, value)
	}
	return result
}

func toOKOutput() *okOutput {
	return &okOutput{Body: OK{OK: true}}
}

type emptyRequest struct{}

type OK struct {
	OK bool `json:"ok"`
}

type okOutput struct {
	Body OK
}
