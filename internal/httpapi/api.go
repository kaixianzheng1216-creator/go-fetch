package httpapi

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/google/uuid"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
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

// DataStore is the persistence contract required by the HTTP API.
type DataStore interface {
	service.UserRepository
	service.CollectionRepository
	service.WebsiteRepository
	service.AnalyticsRepository

	GetUserByID(ctx context.Context, userID uuid.UUID) (domain.User, error)
}

type server struct {
	store    DataStore
	sessions sessionStore
	auth     service.Auth
	collect  service.Collector
	websites service.Websites
	stats    service.Stats
	config   Config
}

func New(dataStore DataStore, sessions *scs.SessionManager, config Config) http.Handler {
	apiServer := server{
		store:    dataStore,
		sessions: sessions,
		auth:     service.NewAuth(dataStore),
		collect:  service.NewCollector(dataStore),
		websites: service.NewWebsites(dataStore),
		stats:    service.NewStats(dataStore),
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
		if _, err := responseWriter.Write(script); err != nil {
			slog.Debug("write tracker script", "error", err)
		}
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
		if err := json.NewEncoder(responseWriter).Encode(huma.ErrorModel{
			Title:  "接口不存在",
			Status: http.StatusNotFound,
			Detail: "接口不存在",
		}); err != nil {
			slog.Debug("write API not found response", "error", err)
		}
		return
	}

	indexHTML, err := webassets.IndexHTML()
	if err != nil {
		http.Error(responseWriter, "dashboard build is missing", http.StatusInternalServerError)
		return
	}

	responseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := responseWriter.Write(indexHTML); err != nil {
		slog.Debug("write dashboard index", "error", err)
	}
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
	return &okOutput{Body: OKResponse{OK: true}}
}

type emptyInput struct{}

type OKResponse struct {
	OK bool `json:"ok"`
}

type okOutput struct {
	Body OKResponse
}
