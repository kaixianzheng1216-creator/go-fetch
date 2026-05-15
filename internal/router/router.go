package router

import (
	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/handler"
)

type Handlers struct {
	Auth    handler.AuthHandler
	Collect handler.CollectHandler
	Website handler.WebsiteHandler
	Stats   handler.StatsHandler
}

func Register(api huma.API, handlers Handlers, authMiddleware huma.Middlewares) {
	RegisterCollect(api, handlers.Collect)
	RegisterAuth(api, handlers.Auth, authMiddleware)
	RegisterWebsite(api, handlers.Website, authMiddleware)
	RegisterStats(api, handlers.Stats, authMiddleware)
}
