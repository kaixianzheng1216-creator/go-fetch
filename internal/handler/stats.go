package handler

import (
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/model"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/service"
)

type StatsHandler struct {
	stats              service.Stats
	currentUser        func(context.Context) model.User
	websiteLookupError func(error) error
}

func NewStats(stats service.Stats, currentUser func(context.Context) model.User, websiteLookupError func(error) error) StatsHandler {
	return StatsHandler{stats: stats, currentUser: currentUser, websiteLookupError: websiteLookupError}
}

type statsRequest struct {
	WebsiteID string `path:"websiteID" format:"uuid"`
	StartAt   int64  `query:"startAt"`
	EndAt     int64  `query:"endAt"`
}

type pageviewsRequest struct {
	WebsiteID string        `path:"websiteID" format:"uuid"`
	StartAt   int64         `query:"startAt"`
	EndAt     int64         `query:"endAt"`
	Unit      DateUnitParam `query:"unit"`
}

type metricsRequest struct {
	WebsiteID string          `path:"websiteID" format:"uuid"`
	StartAt   int64           `query:"startAt"`
	EndAt     int64           `query:"endAt"`
	Type      MetricTypeParam `query:"type" required:"true"`
	Limit     MetricLimit     `query:"limit"`
}

func (handler StatsHandler) GetWebsiteStats(ctx context.Context, input *statsRequest) (*StatsOutput, error) {
	stats, err := handler.stats.WebsiteStats(ctx, handler.currentUser(ctx).ID, input.WebsiteID, OptionalTimeParam(input.StartAt), OptionalTimeParam(input.EndAt))
	if err != nil {
		return nil, handler.statsError(err, "加载统计数据失败")
	}

	return NewStatsOutput(ToWebsiteStats(stats)), nil
}

func (handler StatsHandler) GetWebsitePageviews(ctx context.Context, input *pageviewsRequest) (*PageviewsOutput, error) {
	points, err := handler.stats.WebsitePageviews(ctx, handler.currentUser(ctx).ID, input.WebsiteID, OptionalTimeParam(input.StartAt), OptionalTimeParam(input.EndAt), string(input.Unit))
	if err != nil {
		return nil, handler.statsError(err, "加载页面浏览量失败")
	}

	return NewPageviewsOutput(ToPageviewPoints(points)), nil
}

func (handler StatsHandler) GetWebsiteMetrics(ctx context.Context, input *metricsRequest) (*MetricsOutput, error) {
	rows, err := handler.stats.WebsiteMetrics(ctx, handler.currentUser(ctx).ID, input.WebsiteID, OptionalTimeParam(input.StartAt), OptionalTimeParam(input.EndAt), string(input.Type), int(input.Limit))
	if err != nil {
		if errors.Is(err, model.ErrUnsupportedMetricType) {
			return nil, huma.Error400BadRequest(err.Error())
		}
		return nil, handler.statsError(err, "加载指标数据失败")
	}

	return NewMetricsOutput(ToMetricRows(rows)), nil
}

func (handler StatsHandler) statsError(err error, fallbackMessage string) error {
	if service.IsWebsiteAccessError(err) {
		return handler.websiteLookupError(err)
	}
	return huma.Error500InternalServerError(fallbackMessage)
}
