package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/handler"
)

func registerWebsiteRoutes(api huma.API, websiteHandler handler.WebsiteHandler, authMiddleware huma.Middlewares) {
	huma.Register(
		api,
		requireAuth(
			operation(
				http.MethodGet,
				"/api/websites",
				"listWebsites",
				"列出站点",
				"Websites",
			),
			authMiddleware,
		),
		websiteHandler.ListWebsites,
	)

	createOperation := requireAuth(
		operation(
			http.MethodPost,
			"/api/websites",
			"createWebsite",
			"创建站点",
			"Websites",
		),
		authMiddleware,
	)
	createOperation.DefaultStatus = http.StatusCreated
	huma.Register(api, createOperation, websiteHandler.CreateWebsite)

	huma.Register(
		api,
		requireAuth(
			operation(
				http.MethodGet,
				"/api/websites/{websiteID}",
				"getWebsite",
				"获取站点",
				"Websites",
			),
			authMiddleware,
		),
		websiteHandler.GetWebsite,
	)

	huma.Register(
		api,
		requireAuth(
			operation(
				http.MethodPatch,
				"/api/websites/{websiteID}",
				"updateWebsite",
				"更新站点",
				"Websites",
			),
			authMiddleware,
		),
		websiteHandler.UpdateWebsite,
	)

	huma.Register(
		api,
		requireAuth(
			operation(
				http.MethodDelete,
				"/api/websites/{websiteID}",
				"deleteWebsite",
				"删除站点",
				"Websites",
			),
			authMiddleware,
		),
		websiteHandler.DeleteWebsite,
	)
}
