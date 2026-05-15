package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/handler"
)

func registerStatsRoutes(api huma.API, statsHandler handler.StatsHandler, authMiddleware huma.Middlewares) {
	huma.Register(
		api,
		huma.Operation{
			Method:      http.MethodGet,
			Path:        "/api/websites/{websiteID}/stats",
			OperationID: "websiteStats",
			Summary:     "获取站点统计",
			Tags:        []string{"Analytics"},
			Security:    []map[string][]string{{"sessionCookie": {}}},
			Middlewares: authMiddleware,
		},
		statsHandler.GetWebsiteStats,
	)

	huma.Register(
		api,
		huma.Operation{
			Method:      http.MethodGet,
			Path:        "/api/websites/{websiteID}/pageviews",
			OperationID: "websitePageviews",
			Summary:     "获取页面浏览趋势",
			Tags:        []string{"Analytics"},
			Security:    []map[string][]string{{"sessionCookie": {}}},
			Middlewares: authMiddleware,
		},
		statsHandler.GetWebsitePageviews,
	)

	huma.Register(
		api,
		huma.Operation{
			Method:      http.MethodGet,
			Path:        "/api/websites/{websiteID}/metrics",
			OperationID: "websiteMetrics",
			Summary:     "获取站点指标",
			Tags:        []string{"Analytics"},
			Security:    []map[string][]string{{"sessionCookie": {}}},
			Middlewares: authMiddleware,
		},
		statsHandler.GetWebsiteMetrics,
	)
}
