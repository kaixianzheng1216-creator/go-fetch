package router

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/handler"
)

func registerWebsiteRoutes(api huma.API, websiteHandler handler.WebsiteHandler, authMiddleware huma.Middlewares) {
	huma.Register(
		api,
		huma.Operation{
			Method:      http.MethodGet,
			Path:        "/api/websites",
			OperationID: "listWebsites",
			Summary:     "列出站点",
			Tags:        []string{"Websites"},
			Security:    []map[string][]string{{"sessionCookie": {}}},
			Middlewares: authMiddleware,
		},
		websiteHandler.ListWebsites,
	)

	createOperation := huma.Operation{
		Method:      http.MethodPost,
		Path:        "/api/websites",
		OperationID: "createWebsite",
		Summary:     "创建站点",
		Tags:        []string{"Websites"},
		Security:    []map[string][]string{{"sessionCookie": {}}},
		Middlewares: authMiddleware,
	}
	createOperation.DefaultStatus = http.StatusCreated
	huma.Register(api, createOperation, websiteHandler.CreateWebsite)

	huma.Register(
		api,
		huma.Operation{
			Method:      http.MethodGet,
			Path:        "/api/websites/{websiteID}",
			OperationID: "getWebsite",
			Summary:     "获取站点",
			Tags:        []string{"Websites"},
			Security:    []map[string][]string{{"sessionCookie": {}}},
			Middlewares: authMiddleware,
		},
		websiteHandler.GetWebsite,
	)

	huma.Register(
		api,
		huma.Operation{
			Method:      http.MethodPatch,
			Path:        "/api/websites/{websiteID}",
			OperationID: "updateWebsite",
			Summary:     "更新站点",
			Tags:        []string{"Websites"},
			Security:    []map[string][]string{{"sessionCookie": {}}},
			Middlewares: authMiddleware,
		},
		websiteHandler.UpdateWebsite,
	)

	huma.Register(
		api,
		huma.Operation{
			Method:      http.MethodDelete,
			Path:        "/api/websites/{websiteID}",
			OperationID: "deleteWebsite",
			Summary:     "删除站点",
			Tags:        []string{"Websites"},
			Security:    []map[string][]string{{"sessionCookie": {}}},
			Middlewares: authMiddleware,
		},
		websiteHandler.DeleteWebsite,
	)
}
