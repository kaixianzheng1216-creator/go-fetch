package httpapi

import (
	"context"
	"encoding/json"
	"errors"
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
	sessions sessionStore
	auth     service.AuthService
	collect  service.CollectionService
	users    service.UserService
	websites service.WebsiteService
	stats    service.StatsService
	config   Config
}

type emptyInput struct{}

type OKResponse struct {
	OK bool `json:"ok"`
}

type okOutput struct {
	Body OKResponse
}

type contextKey string

const (
	userContextKey    contextKey = "user"
	requestContextKey contextKey = "request"
)

const (
	errorMessageAPIEndpointNotFound    = "API endpoint not found"
	errorMessageCurrentUserLoadFailed  = "failed to load current user"
	errorMessageDashboardBuildMissing  = "dashboard build is missing"
	errorMessageEventSaveFailed        = "failed to save event"
	errorMessageInvalidCredentials     = "invalid username or password"
	errorMessageLoginSessionCreate     = "failed to create login session"
	errorMessageLogoutFailed           = "failed to log out"
	errorMessageMetricsLoadFailed      = "failed to load metrics"
	errorMessagePageviewsLoadFailed    = "failed to load pageviews"
	errorMessageRequestReadFailed      = "failed to read request"
	errorMessageStatsLoadFailed        = "failed to load stats"
	errorMessageTrackerScriptMissing   = "tracking script is missing"
	errorMessageUnauthenticated        = "authentication required"
	errorMessageUnsupportedEventType   = "unsupported event type"
	errorMessageUserLoadFailed         = "failed to load user"
	errorMessageWebsiteListLoadFailed  = "failed to load websites"
	errorMessageWebsiteLoadFailed      = "failed to load website"
	errorMessageWebsiteNameCannotEmpty = "website name cannot be empty"
	errorMessageWebsiteNotFound        = "website not found"
	errorMessageWebsiteCreateFailed    = "failed to create website"
	errorMessageWebsiteUpdateFailed    = "failed to update website"
)

func New(services Services, sessions *scs.SessionManager, config Config) http.Handler {
	config = config.withDefaults()
	srv := server{
		sessions: sessions,
		auth:     services.Auth,
		collect:  services.Collection,
		users:    services.Users,
		websites: services.Websites,
		stats:    services.Stats,
		config:   config,
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

func (srv server) registerAssets(chiRouter chi.Router) {
	chiRouter.Get("/assets/*", http.FileServer(http.FS(webassets.DashboardFS())).ServeHTTP)
	chiRouter.Get("/script.js", func(responseWriter http.ResponseWriter, _ *http.Request) {
		script, err := webassets.TrackerScript()
		if err != nil {
			http.Error(responseWriter, errorMessageTrackerScriptMissing, http.StatusInternalServerError)
			return
		}

		responseWriter.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		if _, err := responseWriter.Write(script); err != nil {
			slog.Debug("write tracker script", "error", err)
		}
	})
	chiRouter.Get("/*", spaHandler)
}

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

func spaHandler(responseWriter http.ResponseWriter, request *http.Request) {
	switch {
	case request.Method != http.MethodGet:
		http.NotFound(responseWriter, request)
		return
	case strings.HasPrefix(request.URL.Path, "/api/"):
		responseWriter.Header().Set("Content-Type", "application/problem+json")
		responseWriter.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(responseWriter).Encode(huma.ErrorModel{
			Title:  errorMessageAPIEndpointNotFound,
			Status: http.StatusNotFound,
			Detail: errorMessageAPIEndpointNotFound,
		}); err != nil {
			slog.Debug("write API not found response", "error", err)
		}
		return
	}

	indexHTML, err := webassets.IndexHTML()
	if err != nil {
		http.Error(responseWriter, errorMessageDashboardBuildMissing, http.StatusInternalServerError)
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

func publicOperation(method, path, operationID, summary, tag string) huma.Operation {
	return huma.Operation{
		Method:      method,
		Path:        path,
		OperationID: operationID,
		Summary:     summary,
		Tags:        []string{tag},
	}
}

func securedOperation(method, path, operationID, summary, tag string, authMiddleware huma.Middlewares) huma.Operation {
	operation := publicOperation(method, path, operationID, summary, tag)
	operation.Security = []map[string][]string{{"sessionCookie": {}}}
	operation.Middlewares = authMiddleware
	return operation
}

func enumValues(values []string) []any {
	result := make([]any, len(values))
	for i, value := range values {
		result[i] = value
	}
	return result
}

func newOKOutput() *okOutput {
	return &okOutput{Body: OKResponse{OK: true}}
}

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

func OpenAPIJSON() ([]byte, error) {
	chiRouter := chi.NewRouter()
	humaAPI := humachi.New(chiRouter, humaConfig())
	server{config: Config{}.withDefaults()}.registerRoutes(humaAPI)

	return json.MarshalIndent(humaAPI.OpenAPI(), "", "  ")
}

func isNotFound(err error) bool {
	return errors.Is(err, domain.ErrNotFound)
}

func websiteLookupError(err error) error {
	if err == nil {
		return nil
	}
	if isNotFound(err) {
		return huma.Error404NotFound(errorMessageWebsiteNotFound)
	}
	return huma.Error500InternalServerError(errorMessageWebsiteLoadFailed)
}

func websiteMutationError(err error, fallbackMessage string) error {
	if errors.Is(err, service.ErrInvalidWebsiteName) {
		return huma.Error400BadRequest(errorMessageWebsiteNameCannotEmpty)
	}
	return websiteLookupErrorWithFallback(err, fallbackMessage)
}

func websiteLookupErrorWithFallback(err error, fallbackMessage string) error {
	if err == nil {
		return nil
	}
	if isNotFound(err) {
		return huma.Error404NotFound(errorMessageWebsiteNotFound)
	}
	return huma.Error500InternalServerError(fallbackMessage)
}

func collectionError(err error) error {
	switch {
	case errors.Is(err, service.ErrUnsupportedEventType):
		return huma.Error400BadRequest(errorMessageUnsupportedEventType)
	case errors.Is(err, service.ErrMissingClientInfo):
		return huma.Error500InternalServerError(errorMessageRequestReadFailed)
	case isNotFound(err):
		return huma.Error400BadRequest(errorMessageWebsiteNotFound)
	default:
		return huma.Error500InternalServerError(errorMessageEventSaveFailed)
	}
}

func loginError(err error) error {
	if errors.Is(err, service.ErrInvalidCredentials) {
		return huma.Error401Unauthorized(errorMessageInvalidCredentials)
	}
	return huma.Error500InternalServerError(errorMessageUserLoadFailed)
}

func statsError(err error, fallbackMessage string) error {
	if errors.Is(err, domain.ErrUnsupportedDateUnit) || errors.Is(err, domain.ErrUnsupportedMetricType) || errors.Is(err, service.ErrInvalidDateRange) {
		return huma.Error400BadRequest(err.Error())
	}
	if service.IsWebsiteAccessError(err) {
		return websiteLookupError(err)
	}
	return huma.Error500InternalServerError(fallbackMessage)
}
