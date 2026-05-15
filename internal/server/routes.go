package server

import (
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/handler"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/router"
	servermiddleware "github.com/kaixianzheng1216-creator/go-fetch/internal/server/middleware"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/service"
)

func (server *Server) Routes() http.Handler {
	chiRouter := chi.NewRouter()

	chiRouter.Use(chimiddleware.RealIP)
	chiRouter.Use(chimiddleware.RequestID)
	chiRouter.Use(chimiddleware.Recoverer)
	chiRouter.Use(chimiddleware.Logger)
	chiRouter.Use(chimiddleware.Timeout(60 * time.Second))
	chiRouter.Use(server.sessions.LoadAndSave)

	api := humachi.New(chiRouter, humaConfig())
	registerAPIRoutes(api, server)

	chiRouter.Get("/assets/*", server.handleFrontendAsset)
	chiRouter.Get("/script.js", server.handleScript)
	chiRouter.Get("/*", server.handleFrontend)

	return chiRouter
}

func registerAPIRoutes(api huma.API, server *Server) {
	api.UseMiddleware(servermiddleware.CaptureRequest(withRequest))

	authMiddleware := huma.Middlewares{servermiddleware.RequireAuth(api, server.currentUser, withUser)}
	router.Register(api, router.Handlers{
		Auth:    handler.NewAuth(service.NewAuth(server.store, isNotFound), server.sessions, userIDSessionKey, userFromContext),
		Collect: handler.NewCollect(service.NewCollect(server.store), requestFromContext, isNotFound),
		Website: handler.NewWebsite(
			service.NewWebsite(server.store),
			userFromContext,
			websiteLookupError,
		),
		Stats: handler.NewStats(
			service.NewStats(server.store),
			userFromContext,
			websiteLookupError,
		),
	}, authMiddleware)
}
