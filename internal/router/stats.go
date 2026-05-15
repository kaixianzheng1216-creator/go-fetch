package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/handler"
)

func registerStatsRoutes(api huma.API, statsHandler handler.StatsHandler, authMiddleware huma.Middlewares) {
	huma.Register(
		api,
		requireAuth(
			operation(
				http.MethodGet,
				"/api/websites/{websiteID}/stats",
				"websiteStats",
				"获取站点统计",
				"Analytics",
			),
			authMiddleware,
		),
		statsHandler.GetWebsiteStats,
	)

	huma.Register(
		api,
		requireAuth(
			operation(
				http.MethodGet,
				"/api/websites/{websiteID}/pageviews",
				"websitePageviews",
				"获取页面浏览趋势",
				"Analytics",
			),
			authMiddleware,
		),
		statsHandler.GetWebsitePageviews,
	)

	huma.Register(
		api,
		requireAuth(
			operation(
				http.MethodGet,
				"/api/websites/{websiteID}/metrics",
				"websiteMetrics",
				"获取站点指标",
				"Analytics",
			),
			authMiddleware,
		),
		statsHandler.GetWebsiteMetrics,
	)
}
