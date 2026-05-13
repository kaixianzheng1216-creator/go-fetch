package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/httpapi"

	"github.com/danielgtaylor/huma/v2"
)

func registerAnalyticsRoutes(api huma.API, app *App, auth huma.Middlewares) {
	huma.Register(api, authenticated(operation(http.MethodGet, "/api/websites/{websiteID}/stats", "websiteStats", "Analytics", http.StatusUnauthorized, http.StatusNotFound, http.StatusInternalServerError), auth), app.websiteStats)

	huma.Register(api, authenticated(operation(http.MethodGet, "/api/websites/{websiteID}/pageviews", "websitePageviews", "Analytics", http.StatusUnauthorized, http.StatusNotFound, http.StatusInternalServerError), auth), app.websitePageviews)

	huma.Register(api, authenticated(operation(http.MethodGet, "/api/websites/{websiteID}/metrics", "websiteMetrics", "Analytics", http.StatusBadRequest, http.StatusUnauthorized, http.StatusNotFound, http.StatusInternalServerError), auth), app.websiteMetrics)
}

func (a *App) websiteStats(ctx context.Context, input *dateRangeInput) (*jsonBody[httpapi.WebsiteStats], error) {
	if err := a.requireOwnedWebsite(ctx, input.WebsiteID); err != nil {
		return nil, err
	}
	start, end, _ := domain.DateRange(optionalValuePtr(input.StartAt), optionalValuePtr(input.EndAt), "")
	stats, err := a.store.WebsiteStats(ctx, input.WebsiteID, start, end)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to load stats")
	}
	return &jsonBody[httpapi.WebsiteStats]{Body: httpapi.WebsiteStatsFromDomain(stats)}, nil
}

func (a *App) websitePageviews(ctx context.Context, input *pageviewsInput) (*jsonBody[[]httpapi.PageviewPoint], error) {
	if err := a.requireOwnedWebsite(ctx, input.WebsiteID); err != nil {
		return nil, err
	}
	start, end, unit := domain.DateRange(optionalValuePtr(input.StartAt), optionalValuePtr(input.EndAt), string(input.Unit))
	points, err := a.store.Pageviews(ctx, input.WebsiteID, start, end, unit)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to load pageviews")
	}
	return &jsonBody[[]httpapi.PageviewPoint]{Body: httpapi.PageviewPointsFromDomain(points)}, nil
}

func (a *App) websiteMetrics(ctx context.Context, input *metricsInput) (*jsonBody[[]httpapi.MetricRow], error) {
	if err := a.requireOwnedWebsite(ctx, input.WebsiteID); err != nil {
		return nil, err
	}
	start, end, _ := domain.DateRange(optionalValuePtr(input.StartAt), optionalValuePtr(input.EndAt), "")
	metric, ok := domain.ParseMetricType(string(input.Type))
	if !ok {
		return nil, huma.Error400BadRequest(domain.ErrUnsupportedMetricType.Error())
	}
	rows, err := a.store.Metrics(ctx, input.WebsiteID, start, end, metric, int(input.Limit))
	if err != nil {
		if errors.Is(err, domain.ErrUnsupportedMetricType) {
			return nil, huma.Error400BadRequest(err.Error())
		}
		return nil, huma.Error500InternalServerError("failed to load metrics")
	}
	return &jsonBody[[]httpapi.MetricRow]{Body: httpapi.MetricRowsFromDomain(rows)}, nil
}

func (a *App) requireOwnedWebsite(ctx context.Context, websiteID string) error {
	if _, err := a.store.GetWebsite(ctx, userFromContext(ctx).ID, websiteID); err != nil {
		return websiteLookupError(err)
	}
	return nil
}
