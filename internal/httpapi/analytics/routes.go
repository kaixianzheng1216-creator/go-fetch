package analytics

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/httpapi"
)

func Register(api huma.API, h Handler, authMiddleware huma.Middlewares) {
	statsOp := httpapi.NewOperation(
		http.MethodGet,
		"/api/websites/{websiteID}/stats",
		"websiteStats",
		"Analytics",
		http.StatusUnauthorized,
		http.StatusNotFound,
		http.StatusInternalServerError,
	)

	huma.Register(api, httpapi.WithAuth(statsOp, authMiddleware), h.Stats)

	pageviewsOp := httpapi.NewOperation(
		http.MethodGet,
		"/api/websites/{websiteID}/pageviews",
		"websitePageviews",
		"Analytics",
		http.StatusUnauthorized,
		http.StatusNotFound,
		http.StatusInternalServerError,
	)

	huma.Register(api, httpapi.WithAuth(pageviewsOp, authMiddleware), h.Pageviews)

	metricsOp := httpapi.NewOperation(
		http.MethodGet,
		"/api/websites/{websiteID}/metrics",
		"websiteMetrics",
		"Analytics",
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusNotFound,
		http.StatusInternalServerError,
	)

	huma.Register(api, httpapi.WithAuth(metricsOp, authMiddleware), h.Metrics)
}
