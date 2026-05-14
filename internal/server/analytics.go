package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"

	"github.com/danielgtaylor/huma/v2"
)

func registerAnalyticsRoutes(api huma.API, app *App, auth huma.Middlewares) {
	statsOp := operation(
		http.MethodGet,
		"/api/websites/{websiteID}/stats",
		"websiteStats",
		"Analytics",
		http.StatusUnauthorized,
		http.StatusNotFound,
		http.StatusInternalServerError,
	)

	huma.Register(api, authenticated(statsOp, auth), app.websiteStats)

	pageviewsOp := operation(
		http.MethodGet,
		"/api/websites/{websiteID}/pageviews",
		"websitePageviews",
		"Analytics",
		http.StatusUnauthorized,
		http.StatusNotFound,
		http.StatusInternalServerError,
	)

	huma.Register(api, authenticated(pageviewsOp, auth), app.websitePageviews)

	metricsOp := operation(
		http.MethodGet,
		"/api/websites/{websiteID}/metrics",
		"websiteMetrics",
		"Analytics",
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusNotFound,
		http.StatusInternalServerError,
	)

	huma.Register(api, authenticated(metricsOp, auth), app.websiteMetrics)
}

func (a *App) websiteStats(ctx context.Context, input *dateRangeInput) (*jsonBody[WebsiteStats], error) {
	if err := a.requireOwnedWebsite(ctx, input.WebsiteID); err != nil {
		return nil, err
	}

	start, end, _ := domain.DateRange(queryTimePtr(input.StartAt), queryTimePtr(input.EndAt), "")
	stats, err := a.store.WebsiteStats(ctx, input.WebsiteID, start, end)
	if err != nil {
		return nil, huma.Error500InternalServerError("加载统计数据失败")
	}

	response := WebsiteStatsFromDomain(stats)

	return jsonResponse(response), nil
}

func (a *App) websitePageviews(ctx context.Context, input *pageviewsInput) (*jsonBody[[]PageviewPoint], error) {
	if err := a.requireOwnedWebsite(ctx, input.WebsiteID); err != nil {
		return nil, err
	}

	start, end, unit := domain.DateRange(queryTimePtr(input.StartAt), queryTimePtr(input.EndAt), string(input.Unit))
	points, err := a.store.Pageviews(ctx, input.WebsiteID, start, end, unit)
	if err != nil {
		return nil, huma.Error500InternalServerError("加载浏览量数据失败")
	}

	response := PageviewPointsFromDomain(points)

	return jsonResponse(response), nil
}

func (a *App) websiteMetrics(ctx context.Context, input *metricsInput) (*jsonBody[[]MetricRow], error) {
	if err := a.requireOwnedWebsite(ctx, input.WebsiteID); err != nil {
		return nil, err
	}

	start, end, _ := domain.DateRange(queryTimePtr(input.StartAt), queryTimePtr(input.EndAt), "")
	metric, ok := domain.ParseMetricType(string(input.Type))
	if !ok {
		return nil, huma.Error400BadRequest(domain.ErrUnsupportedMetricType.Error())
	}

	limit := int(input.Limit)
	if limit == 0 {
		limit = domain.DefaultMetricLimit
	}

	rows, err := a.store.Metrics(ctx, input.WebsiteID, start, end, metric, limit)
	if err != nil {
		if errors.Is(err, domain.ErrUnsupportedMetricType) {
			return nil, huma.Error400BadRequest(err.Error())
		}

		return nil, huma.Error500InternalServerError("加载指标数据失败")
	}

	response := MetricRowsFromDomain(rows)

	return jsonResponse(response), nil
}

func (a *App) requireOwnedWebsite(ctx context.Context, websiteID string) error {
	if _, err := a.store.GetWebsite(ctx, userFromContext(ctx).ID, websiteID); err != nil {
		return websiteLookupError(err)
	}
	return nil
}
