package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/handler"
)

func RegisterWebsite(api huma.API, websiteHandler handler.WebsiteHandler, authMiddleware huma.Middlewares) {
	huma.Register(
		api,
		WithAuth(
			NewOperation(
				http.MethodGet,
				"/api/websites",
				"listWebsites",
				"列出站点",
				"Websites",
				http.StatusUnauthorized,
				http.StatusInternalServerError,
			),
			authMiddleware,
		),
		websiteHandler.ListWebsites,
	)

	huma.Register(
		api,
		WithDefaultStatus(
			WithAuth(
				NewOperation(
					http.MethodPost,
					"/api/websites",
					"createWebsite",
					"创建站点",
					"Websites",
					http.StatusBadRequest,
					http.StatusUnauthorized,
					http.StatusUnprocessableEntity,
					http.StatusInternalServerError,
				),
				authMiddleware,
			),
			http.StatusCreated,
		),
		websiteHandler.CreateWebsite,
	)

	huma.Register(
		api,
		WithAuth(
			NewOperation(
				http.MethodGet,
				"/api/websites/{websiteID}",
				"getWebsite",
				"获取站点",
				"Websites",
				http.StatusUnauthorized,
				http.StatusNotFound,
				http.StatusInternalServerError,
			),
			authMiddleware,
		),
		websiteHandler.GetWebsite,
	)

	huma.Register(
		api,
		WithAuth(
			NewOperation(
				http.MethodPatch,
				"/api/websites/{websiteID}",
				"updateWebsite",
				"更新站点",
				"Websites",
				http.StatusBadRequest,
				http.StatusUnauthorized,
				http.StatusNotFound,
				http.StatusUnprocessableEntity,
				http.StatusInternalServerError,
			),
			authMiddleware,
		),
		websiteHandler.UpdateWebsite,
	)

	huma.Register(
		api,
		WithAuth(
			NewOperation(
				http.MethodDelete,
				"/api/websites/{websiteID}",
				"deleteWebsite",
				"删除站点",
				"Websites",
				http.StatusUnauthorized,
				http.StatusNotFound,
				http.StatusInternalServerError,
			),
			authMiddleware,
		),
		websiteHandler.DeleteWebsite,
	)
}
