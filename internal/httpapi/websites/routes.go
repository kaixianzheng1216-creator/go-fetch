package websites

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/httpapi"
)

func Register(api huma.API, handler Handler, authMiddleware huma.Middlewares) {
	listOp := httpapi.NewOperation(
		http.MethodGet,
		"/api/websites",
		"listWebsites",
		"Websites",
		http.StatusUnauthorized,
		http.StatusInternalServerError,
	)

	huma.Register(api, httpapi.WithAuth(listOp, authMiddleware), handler.ListWebsites)

	createOp := httpapi.NewOperation(
		http.MethodPost,
		"/api/websites",
		"createWebsite",
		"Websites",
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusUnprocessableEntity,
		http.StatusInternalServerError,
	)

	createOp = httpapi.WithAuth(createOp, authMiddleware)
	createOp.DefaultStatus = http.StatusCreated

	huma.Register(api, createOp, handler.CreateWebsite)

	getOp := httpapi.NewOperation(
		http.MethodGet,
		"/api/websites/{websiteID}",
		"getWebsite",
		"Websites",
		http.StatusUnauthorized,
		http.StatusNotFound,
		http.StatusInternalServerError,
	)

	huma.Register(api, httpapi.WithAuth(getOp, authMiddleware), handler.GetWebsite)

	updateOp := httpapi.NewOperation(
		http.MethodPatch,
		"/api/websites/{websiteID}",
		"updateWebsite",
		"Websites",
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusNotFound,
		http.StatusUnprocessableEntity,
		http.StatusInternalServerError,
	)

	huma.Register(api, httpapi.WithAuth(updateOp, authMiddleware), handler.UpdateWebsite)

	deleteOp := httpapi.NewOperation(
		http.MethodDelete,
		"/api/websites/{websiteID}",
		"deleteWebsite",
		"Websites",
		http.StatusUnauthorized,
		http.StatusNotFound,
		http.StatusInternalServerError,
	)

	huma.Register(api, httpapi.WithAuth(deleteOp, authMiddleware), handler.DeleteWebsite)
}
