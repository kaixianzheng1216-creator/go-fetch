package handler

import (
	"context"
	"errors"

	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	requestdto "github.com/kaixianzheng1216-creator/go-fetch/internal/dto/request"
	responsedto "github.com/kaixianzheng1216-creator/go-fetch/internal/dto/response"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/service"
)

type StatsHandler struct {
	stats              service.Stats
	currentUser        func(context.Context) domain.User
	websiteLookupError func(error) error
}

func NewStats(stats service.Stats, currentUser func(context.Context) domain.User, websiteLookupError func(error) error) StatsHandler {
	return StatsHandler{stats: stats, currentUser: currentUser, websiteLookupError: websiteLookupError}
}

type statsRequest struct {
	WebsiteID string `path:"websiteID" format:"uuid"`
	StartAt   int64  `query:"startAt"`
	EndAt     int64  `query:"endAt"`
}

type pageviewsRequest struct {
	WebsiteID string                   `path:"websiteID" format:"uuid"`
	StartAt   int64                    `query:"startAt"`
	EndAt     int64                    `query:"endAt"`
	Unit      requestdto.DateUnitParam `query:"unit"`
}

type metricsRequest struct {
	WebsiteID string                     `path:"websiteID" format:"uuid"`
	StartAt   int64                      `query:"startAt"`
	EndAt     int64                      `query:"endAt"`
	Type      requestdto.MetricTypeParam `query:"type" required:"true"`
	Limit     requestdto.MetricLimit     `query:"limit"`
}

func (handler StatsHandler) GetWebsiteStats(ctx context.Context, input *statsRequest) (*responsedto.StatsOutput, error) {
	stats, err := handler.stats.WebsiteStats(ctx, handler.currentUser(ctx).ID, input.WebsiteID, requestdto.OptionalTimeParam(input.StartAt), requestdto.OptionalTimeParam(input.EndAt))
	if err != nil {
		return nil, handler.statsError(err, "加载统计数据失败")
	}

	return responsedto.NewStatsOutput(responsedto.ToWebsiteStats(stats)), nil
}

func (handler StatsHandler) GetWebsitePageviews(ctx context.Context, input *pageviewsRequest) (*responsedto.PageviewsOutput, error) {
	points, err := handler.stats.WebsitePageviews(ctx, handler.currentUser(ctx).ID, input.WebsiteID, requestdto.OptionalTimeParam(input.StartAt), requestdto.OptionalTimeParam(input.EndAt), string(input.Unit))
	if err != nil {
		return nil, handler.statsError(err, "加载页面浏览量失败")
	}

	return responsedto.NewPageviewsOutput(responsedto.ToPageviewPoints(points)), nil
}

func (handler StatsHandler) GetWebsiteMetrics(ctx context.Context, input *metricsRequest) (*responsedto.MetricsOutput, error) {
	rows, err := handler.stats.WebsiteMetrics(ctx, handler.currentUser(ctx).ID, input.WebsiteID, requestdto.OptionalTimeParam(input.StartAt), requestdto.OptionalTimeParam(input.EndAt), string(input.Type), int(input.Limit))
	if err != nil {
		if errors.Is(err, domain.ErrUnsupportedMetricType) {
			return nil, huma.Error400BadRequest(err.Error())
		}
		return nil, handler.statsError(err, "加载指标数据失败")
	}

	return responsedto.NewMetricsOutput(responsedto.ToMetricRows(rows)), nil
}

func (handler StatsHandler) statsError(err error, fallbackMessage string) error {
	if service.IsWebsiteAccessError(err) {
		return handler.websiteLookupError(err)
	}
	return huma.Error500InternalServerError(fallbackMessage)
}
