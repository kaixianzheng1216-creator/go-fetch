package httpapi

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/service"
)

type getWebsiteStatsInput struct {
	WebsiteID uuid.UUID `path:"websiteID" format:"uuid"`
	StartAt   int64     `query:"startAt"`
	EndAt     int64     `query:"endAt"`
}

type getWebsitePageviewsInput struct {
	WebsiteID uuid.UUID     `path:"websiteID" format:"uuid"`
	StartAt   int64         `query:"startAt"`
	EndAt     int64         `query:"endAt"`
	Unit      dateUnitParam `query:"unit"`
}

type getWebsiteMetricsInput struct {
	WebsiteID uuid.UUID        `path:"websiteID" format:"uuid"`
	StartAt   int64            `query:"startAt"`
	EndAt     int64            `query:"endAt"`
	Type      metricTypeParam  `query:"type" required:"true"`
	Limit     metricLimitParam `query:"limit"`
}

type dateUnitParam string

func (dateUnitParam) Schema(huma.Registry) *huma.Schema {
	return &huma.Schema{
		Type:    huma.TypeString,
		Enum:    enumValues(domain.DateUnitValues()),
		Default: string(domain.DefaultDateUnit),
	}
}

type metricTypeParam string

func (metricTypeParam) Schema(huma.Registry) *huma.Schema {
	return &huma.Schema{
		Type: huma.TypeString,
		Enum: enumValues(domain.MetricTypeValues()),
	}
}

type metricLimitParam int

func (metricLimitParam) Schema(huma.Registry) *huma.Schema {
	minValue := 1.0
	maxValue := float64(domain.MaxMetricLimit)
	return &huma.Schema{
		Type:    huma.TypeInteger,
		Default: domain.DefaultMetricLimit,
		Minimum: &minValue,
		Maximum: &maxValue,
	}
}

type getWebsiteStatsOutput struct {
	Body struct {
		Pageviews       int64 `json:"pageviews"`
		Visitors        int64 `json:"visitors"`
		Visits          int64 `json:"visits"`
		Bounces         int64 `json:"bounces"`
		TotalTime       int64 `json:"totalTime"`
		AvgVisitSeconds int64 `json:"avgVisitSeconds"`
	}
}

type getWebsitePageviewsOutput struct {
	Body pageviewListBody
}

type getWebsiteMetricsOutput struct {
	Body metricListBody
}

type pageviewListBody []struct {
	Time     time.Time `json:"time"`
	Label    string    `json:"label"`
	Views    int64     `json:"views"`
	Visitors int64     `json:"visitors"`
}

type metricListBody []struct {
	Name     string `json:"name"`
	Views    int64  `json:"views"`
	Visitors int64  `json:"visitors"`
}

func (srv server) registerStatsRoutes(humaAPI huma.API, authMiddleware huma.Middlewares) {
	huma.Register(
		humaAPI,
		securedOperation(http.MethodGet, "/api/websites/{websiteID}/stats", "websiteStats", "Get website stats", "Analytics", authMiddleware),
		srv.getWebsiteStats,
	)

	huma.Register(
		humaAPI,
		securedOperation(http.MethodGet, "/api/websites/{websiteID}/pageviews", "websitePageviews", "Get website pageviews", "Analytics", authMiddleware),
		srv.getWebsitePageviews,
	)

	huma.Register(
		humaAPI,
		securedOperation(http.MethodGet, "/api/websites/{websiteID}/metrics", "websiteMetrics", "Get website metrics", "Analytics", authMiddleware),
		srv.getWebsiteMetrics,
	)
}

func (srv server) getWebsiteStats(ctx context.Context, input *getWebsiteStatsInput) (*getWebsiteStatsOutput, error) {
	query, err := newStatsQuery(ctx, input.WebsiteID, input.StartAt, input.EndAt)
	if err != nil {
		return nil, err
	}

	stats, err := srv.stats.Summary(ctx, query)
	if err != nil {
		return nil, statsError(err, errorMessageStatsLoadFailed)
	}

	return newWebsiteStatsOutput(stats), nil
}

func (srv server) getWebsitePageviews(ctx context.Context, input *getWebsitePageviewsInput) (*getWebsitePageviewsOutput, error) {
	query, err := newStatsQuery(ctx, input.WebsiteID, input.StartAt, input.EndAt)
	if err != nil {
		return nil, err
	}

	buckets, err := srv.stats.Pageviews(ctx, service.PageviewsQuery{
		StatsQuery: query,
		Unit:       domain.DateUnit(input.Unit),
	})
	if err != nil {
		return nil, statsError(err, errorMessagePageviewsLoadFailed)
	}

	return newWebsitePageviewsOutput(buckets), nil
}

func (srv server) getWebsiteMetrics(ctx context.Context, input *getWebsiteMetricsInput) (*getWebsiteMetricsOutput, error) {
	query, err := newStatsQuery(ctx, input.WebsiteID, input.StartAt, input.EndAt)
	if err != nil {
		return nil, err
	}

	metrics, err := srv.stats.Metrics(ctx, service.MetricsQuery{
		StatsQuery: query,
		Type:       domain.MetricType(input.Type),
		Limit:      int(input.Limit),
	})
	if err != nil {
		return nil, statsError(err, errorMessageMetricsLoadFailed)
	}

	return newWebsiteMetricsOutput(metrics), nil
}

func newStatsQuery(ctx context.Context, websiteID uuid.UUID, startAt, endAt int64) (service.StatsQuery, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return service.StatsQuery{}, err
	}

	return service.StatsQuery{
		UserID:    userID,
		WebsiteID: websiteID,
		Range:     dateRangeFromInput(startAt, endAt),
	}, nil
}

func dateRangeFromInput(startAt, endAt int64) service.DateRange {
	return service.DateRange{
		StartAt: optionalTimeParam(startAt),
		EndAt:   optionalTimeParam(endAt),
	}
}

func optionalTimeParam(value int64) *time.Time {
	if value == 0 {
		return nil
	}
	timestamp := time.UnixMilli(value).UTC()
	return &timestamp
}

func newWebsiteStatsOutput(stats domain.WebsiteStats) *getWebsiteStatsOutput {
	output := &getWebsiteStatsOutput{}
	output.Body.Pageviews = stats.Pageviews
	output.Body.Visitors = stats.Visitors
	output.Body.Visits = stats.Visits
	output.Body.Bounces = stats.Bounces
	output.Body.TotalTime = stats.TotalTime
	output.Body.AvgVisitSeconds = stats.AvgVisitSeconds
	return output
}

func newWebsitePageviewsOutput(buckets []domain.PageviewBucket) *getWebsitePageviewsOutput {
	output := &getWebsitePageviewsOutput{Body: make(pageviewListBody, len(buckets))}
	for i, bucket := range buckets {
		output.Body[i].Time = bucket.Time
		output.Body[i].Label = bucket.Label
		output.Body[i].Views = bucket.Views
		output.Body[i].Visitors = bucket.Visitors
	}
	return output
}

func newWebsiteMetricsOutput(metrics []domain.Metric) *getWebsiteMetricsOutput {
	output := &getWebsiteMetricsOutput{Body: make(metricListBody, len(metrics))}
	for i, metric := range metrics {
		output.Body[i].Name = metric.Name
		output.Body[i].Views = metric.Views
		output.Body[i].Visitors = metric.Visitors
	}
	return output
}
