package analytics

import (
	"context"
	"errors"
	"time"

	"github.com/danielgtaylor/huma/v2"

	eventdomain "github.com/kaixianzheng1216-creator/go-fetch/internal/event"
	userdomain "github.com/kaixianzheng1216-creator/go-fetch/internal/user"
	websitedomain "github.com/kaixianzheng1216-creator/go-fetch/internal/website"
)

type Store interface {
	GetWebsite(ctx context.Context, userID, websiteID string) (websitedomain.Website, error)
	WebsiteStats(ctx context.Context, websiteID string, start, end time.Time) (eventdomain.WebsiteStats, error)
	WebsitePageviews(ctx context.Context, websiteID string, start, end time.Time, unit eventdomain.DateUnit) ([]eventdomain.PageviewPoint, error)
	WebsiteMetrics(ctx context.Context, websiteID string, start, end time.Time, metric eventdomain.MetricType, limit int) ([]eventdomain.MetricRow, error)
}

type Handler struct {
	store              Store
	currentUser        func(context.Context) userdomain.User
	websiteLookupError func(error) error
}

func New(
	dataStore Store,
	currentUser func(context.Context) userdomain.User,
	websiteLookupError func(error) error,
) Handler {
	return Handler{
		store:              dataStore,
		currentUser:        currentUser,
		websiteLookupError: websiteLookupError,
	}
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

func (handler Handler) GetWebsiteStats(ctx context.Context, request *statsRequest) (*statsOutput, error) {
	if err := handler.requireWebsiteAccess(ctx, request.WebsiteID); err != nil {
		return nil, err
	}

	start, end, _ := eventdomain.DateRange(OptionalTimeParam(request.StartAt), OptionalTimeParam(request.EndAt), "")
	stats, err := handler.store.WebsiteStats(ctx, request.WebsiteID, start, end)
	if err != nil {
		return nil, huma.Error500InternalServerError("加载统计数据失败")
	}

	return newStatsOutput(ToWebsiteStats(stats)), nil
}

func (handler Handler) GetWebsitePageviews(ctx context.Context, request *pageviewsRequest) (*pageviewsOutput, error) {
	if err := handler.requireWebsiteAccess(ctx, request.WebsiteID); err != nil {
		return nil, err
	}

	start, end, unit := eventdomain.DateRange(OptionalTimeParam(request.StartAt), OptionalTimeParam(request.EndAt), string(request.Unit))
	points, err := handler.store.WebsitePageviews(ctx, request.WebsiteID, start, end, unit)
	if err != nil {
		return nil, huma.Error500InternalServerError("加载页面浏览量失败")
	}

	return newPageviewsOutput(ToPageviewPoints(points)), nil
}

func (handler Handler) GetWebsiteMetrics(ctx context.Context, request *metricsRequest) (*metricsOutput, error) {
	if err := handler.requireWebsiteAccess(ctx, request.WebsiteID); err != nil {
		return nil, err
	}

	start, end, _ := eventdomain.DateRange(OptionalTimeParam(request.StartAt), OptionalTimeParam(request.EndAt), "")
	metric, isSupportedMetricType := eventdomain.ParseMetricType(string(request.Type))
	if !isSupportedMetricType {
		return nil, huma.Error400BadRequest(eventdomain.ErrUnsupportedMetricType.Error())
	}

	limit := int(request.Limit)
	if limit == 0 {
		limit = eventdomain.DefaultMetricLimit
	}

	rows, err := handler.store.WebsiteMetrics(ctx, request.WebsiteID, start, end, metric, limit)
	if err != nil {
		if errors.Is(err, eventdomain.ErrUnsupportedMetricType) {
			return nil, huma.Error400BadRequest(err.Error())
		}

		return nil, huma.Error500InternalServerError("加载指标数据失败")
	}

	return newMetricsOutput(ToMetricRows(rows)), nil
}

func (handler Handler) requireWebsiteAccess(ctx context.Context, websiteID string) error {
	if _, err := handler.store.GetWebsite(ctx, handler.currentUser(ctx).ID, websiteID); err != nil {
		return handler.websiteLookupError(err)
	}
	return nil
}
