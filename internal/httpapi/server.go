package httpapi

import (
	"context"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/service"
)

type Config struct {
	CollectCORSAllowedOrigins []string
	RequestTimeout            time.Duration
	TrustProxyHeaders         bool
}

type Services struct {
	Auth       service.AuthService
	Collection service.CollectionService
	Stats      service.StatsService
	Users      service.UserService
	Websites   service.WebsiteService
}

type sessionStore interface {
	RenewToken(ctx context.Context) error
	Put(ctx context.Context, key string, val any)
	Destroy(ctx context.Context) error
	GetString(ctx context.Context, key string) string
}

type server struct {
	sessions   sessionStore
	auth       service.AuthService
	collection service.CollectionService
	users      service.UserService
	websites   service.WebsiteService
	stats      service.StatsService
	config     Config
}

type emptyInput struct{}

type okOutput struct {
	Body struct {
		OK bool `json:"ok"`
	}
}

func New(services Services, sessions *scs.SessionManager, config Config) http.Handler {
	config = config.withDefaults()
	srv := server{
		sessions:   sessions,
		auth:       services.Auth,
		collection: services.Collection,
		users:      services.Users,
		websites:   services.Websites,
		stats:      services.Stats,
		config:     config,
	}

	chiRouter := chi.NewRouter()
	if config.TrustProxyHeaders {
		chiRouter.Use(chimiddleware.RealIP)
	}
	chiRouter.Use(chimiddleware.RequestID)
	chiRouter.Use(chimiddleware.Recoverer)
	chiRouter.Use(chimiddleware.Logger)
	chiRouter.Use(chimiddleware.Timeout(config.RequestTimeout))
	chiRouter.Use(srv.collectCORSMiddleware)
	if sessions != nil {
		chiRouter.Use(sessions.LoadAndSave)
	}

	humaAPI := humachi.New(chiRouter, humaConfig())
	humaAPI.UseMiddleware(captureRequest)
	srv.registerRoutes(humaAPI)
	srv.registerAssets(chiRouter)

	return chiRouter
}

func (config Config) withDefaults() Config {
	if len(config.CollectCORSAllowedOrigins) == 0 {
		config.CollectCORSAllowedOrigins = []string{"*"}
	}
	if config.RequestTimeout <= 0 {
		config.RequestTimeout = 60 * time.Second
	}
	return config
}

func (srv server) registerRoutes(humaAPI huma.API) {
	authMiddleware := huma.Middlewares{srv.requireAuth(humaAPI)}

	srv.registerCollectRoutes(humaAPI)
	srv.registerAuthRoutes(humaAPI, authMiddleware)
	srv.registerWebsiteRoutes(humaAPI, authMiddleware)
	srv.registerStatsRoutes(humaAPI, authMiddleware)
}

func newOKOutput() *okOutput {
	output := &okOutput{}
	output.Body.OK = true
	return output
}
