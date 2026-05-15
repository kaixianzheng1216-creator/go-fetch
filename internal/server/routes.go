package server

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/middleware"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/server/handler/eventhandler"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/server/handler/summaryhandler"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/server/handler/userhandler"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/server/handler/websitehandler"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/session"
)

func (a *App) Routes() http.Handler {
	r := chi.NewRouter()

	middleware.UseHTTP(r, a.sessions)

	api := humachi.New(r, humaConfig())
	registerAPIRoutes(api, a)

	r.Get("/assets/*", a.handleFrontendAsset)
	r.Get("/script.js", a.handleScript)
	r.Get("/*", a.handleFrontend)

	return r
}

func registerAPIRoutes(api huma.API, app *App) {
	api.UseMiddleware(middleware.CaptureRequest(withRequest))
	authMiddleware := huma.Middlewares{middleware.RequireAuth(api, app.currentUser, withUser)}

	authHandler := userhandler.New(app.store, app.sessions, session.UserIDKey, userFromContext, isNotFound)
	collectHandler := eventhandler.New(app.store, requestFromContext, isNotFound)
	websitesHandler := websitehandler.New(app.store, userFromContext, websiteLookupError)
	analyticsHandler := summaryhandler.New(app.store, userFromContext, websiteLookupError)

	registerCollectRoutes(api, collectHandler)
	registerAuthRoutes(api, authHandler, authMiddleware)
	registerWebsiteRoutes(api, websitesHandler, authMiddleware)
	registerAnalyticsRoutes(api, analyticsHandler, authMiddleware)
}

func registerCollectRoutes(api huma.API, h eventhandler.Handler) {
	collectOp := newOperation(
		http.MethodPost,
		"/api/collect",
		"collect",
		"Collection",
		http.StatusBadRequest,
		http.StatusUnprocessableEntity,
		http.StatusInternalServerError,
	)

	collectOp.MaxBodyBytes = 256 * 1024

	huma.Register(api, collectOp, h.Collect)
}

func registerAuthRoutes(api huma.API, h userhandler.Handler, authMiddleware huma.Middlewares) {
	loginOp := newOperation(
		http.MethodPost,
		"/api/login",
		"login",
		"Auth",
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusUnprocessableEntity,
		http.StatusInternalServerError,
	)

	huma.Register(api, loginOp, h.Login)

	logoutOp := newOperation(
		http.MethodPost,
		"/api/logout",
		"logout",
		"Auth",
		http.StatusInternalServerError,
	)

	huma.Register(api, logoutOp, h.Logout)

	meOp := newOperation(
		http.MethodGet,
		"/api/me",
		"getCurrentUser",
		"Auth",
		http.StatusUnauthorized,
	)

	huma.Register(api, withAuth(meOp, authMiddleware), h.Me)
}

func registerWebsiteRoutes(api huma.API, h websitehandler.Handler, authMiddleware huma.Middlewares) {
	listOp := newOperation(
		http.MethodGet,
		"/api/websites",
		"listWebsites",
		"Websites",
		http.StatusUnauthorized,
		http.StatusInternalServerError,
	)

	huma.Register(api, withAuth(listOp, authMiddleware), h.List)

	createOp := newOperation(
		http.MethodPost,
		"/api/websites",
		"createWebsite",
		"Websites",
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusUnprocessableEntity,
		http.StatusInternalServerError,
	)

	createOp = withAuth(createOp, authMiddleware)
	createOp.DefaultStatus = http.StatusCreated

	huma.Register(api, createOp, h.Create)

	getOp := newOperation(
		http.MethodGet,
		"/api/websites/{websiteID}",
		"getWebsite",
		"Websites",
		http.StatusUnauthorized,
		http.StatusNotFound,
		http.StatusInternalServerError,
	)

	huma.Register(api, withAuth(getOp, authMiddleware), h.Get)

	updateOp := newOperation(
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

	huma.Register(api, withAuth(updateOp, authMiddleware), h.Update)

	deleteOp := newOperation(
		http.MethodDelete,
		"/api/websites/{websiteID}",
		"deleteWebsite",
		"Websites",
		http.StatusUnauthorized,
		http.StatusNotFound,
		http.StatusInternalServerError,
	)

	huma.Register(api, withAuth(deleteOp, authMiddleware), h.Delete)
}

func registerAnalyticsRoutes(api huma.API, h summaryhandler.Handler, authMiddleware huma.Middlewares) {
	statsOp := newOperation(
		http.MethodGet,
		"/api/websites/{websiteID}/stats",
		"websiteStats",
		"Analytics",
		http.StatusUnauthorized,
		http.StatusNotFound,
		http.StatusInternalServerError,
	)

	huma.Register(api, withAuth(statsOp, authMiddleware), h.Stats)

	pageviewsOp := newOperation(
		http.MethodGet,
		"/api/websites/{websiteID}/pageviews",
		"websitePageviews",
		"Analytics",
		http.StatusUnauthorized,
		http.StatusNotFound,
		http.StatusInternalServerError,
	)

	huma.Register(api, withAuth(pageviewsOp, authMiddleware), h.Pageviews)

	metricsOp := newOperation(
		http.MethodGet,
		"/api/websites/{websiteID}/metrics",
		"websiteMetrics",
		"Analytics",
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusNotFound,
		http.StatusInternalServerError,
	)

	huma.Register(api, withAuth(metricsOp, authMiddleware), h.Metrics)
}

func newOperation(method, path, operationID, tag string, statusCodes ...int) huma.Operation {
	return huma.Operation{
		Method:      method,
		Path:        path,
		OperationID: operationID,
		Tags:        []string{tag},
		Errors:      statusCodes,
	}
}

func withAuth(op huma.Operation, middlewares huma.Middlewares) huma.Operation {
	op.Security = []map[string][]string{{"sessionCookie": {}}}
	op.Middlewares = append(op.Middlewares, middlewares...)
	return op
}
