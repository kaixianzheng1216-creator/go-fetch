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

type websiteStatsInput struct {
	WebsiteID uuid.UUID `path:"websiteID" format:"uuid"`
	StartAt   int64     `query:"startAt"`
	EndAt     int64     `query:"endAt"`
}

type websitePageviewsInput struct {
	WebsiteID uuid.UUID     `path:"websiteID" format:"uuid"`
	StartAt   int64         `query:"startAt"`
	EndAt     int64         `query:"endAt"`
	Unit      DateUnitParam `query:"unit"`
}

type websiteMetricsInput struct {
	WebsiteID uuid.UUID        `path:"websiteID" format:"uuid"`
	StartAt   int64            `query:"startAt"`
	EndAt     int64            `query:"endAt"`
	Type      MetricTypeParam  `query:"type" required:"true"`
	Limit     MetricLimitParam `query:"limit"`
}

type DateUnitParam string

func (DateUnitParam) Schema(huma.Registry) *huma.Schema {
	return &huma.Schema{
		Type:    huma.TypeString,
		Enum:    enumValues(domain.DateUnitValues()),
		Default: string(domain.DefaultDateUnit),
	}
}

type MetricTypeParam string

func (MetricTypeParam) Schema(huma.Registry) *huma.Schema {
	return &huma.Schema{
		Type: huma.TypeString,
		Enum: enumValues(domain.MetricTypeValues()),
	}
}

type MetricLimitParam int

func (MetricLimitParam) Schema(huma.Registry) *huma.Schema {
	minValue := 1.0
	maxValue := float64(domain.MaxMetricLimit)
	return &huma.Schema{
		Type:    huma.TypeInteger,
		Default: domain.DefaultMetricLimit,
		Minimum: &minValue,
		Maximum: &maxValue,
	}
}

type WebsiteStatsResponse struct {
	Pageviews       int64 `json:"pageviews"`
	Visitors        int64 `json:"visitors"`
	Visits          int64 `json:"visits"`
	Bounces         int64 `json:"bounces"`
	TotalTime       int64 `json:"totalTime"`
	AvgVisitSeconds int64 `json:"avgVisitSeconds"`
}

type PageviewResponse struct {
	Time     time.Time `json:"time"`
	Label    string    `json:"label"`
	Views    int64     `json:"views"`
	Visitors int64     `json:"visitors"`
}

type MetricResponse struct {
	Name     string `json:"name"`
	Views    int64  `json:"views"`
	Visitors int64  `json:"visitors"`
}

type websiteStatsOutput struct {
	Body WebsiteStatsResponse
}

type websitePageviewsOutput struct {
	Body []PageviewResponse
}

type websiteMetricsOutput struct {
	Body []MetricResponse
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

func (srv server) getWebsiteStats(ctx context.Context, input *websiteStatsInput) (*websiteStatsOutput, error) {
	stats, err := srv.stats.Summary(ctx, service.StatsQuery{
		UserID:    currentUser(ctx).ID,
		WebsiteID: input.WebsiteID,
		Range:     dateRangeFromInput(input.StartAt, input.EndAt),
	})
	if err != nil {
		return nil, statsError(err, errorMessageStatsLoadFailed)
	}

	return &websiteStatsOutput{Body: newWebsiteStatsResponse(stats)}, nil
}

func (srv server) getWebsitePageviews(ctx context.Context, input *websitePageviewsInput) (*websitePageviewsOutput, error) {
	buckets, err := srv.stats.Pageviews(ctx, service.PageviewsQuery{
		StatsQuery: service.StatsQuery{
			UserID:    currentUser(ctx).ID,
			WebsiteID: input.WebsiteID,
			Range:     dateRangeFromInput(input.StartAt, input.EndAt),
		},
		Unit: domain.DateUnit(input.Unit),
	})
	if err != nil {
		return nil, statsError(err, errorMessagePageviewsLoadFailed)
	}

	return &websitePageviewsOutput{Body: newPageviewResponses(buckets)}, nil
}

func (srv server) getWebsiteMetrics(ctx context.Context, input *websiteMetricsInput) (*websiteMetricsOutput, error) {
	metrics, err := srv.stats.Metrics(ctx, service.MetricsQuery{
		StatsQuery: service.StatsQuery{
			UserID:    currentUser(ctx).ID,
			WebsiteID: input.WebsiteID,
			Range:     dateRangeFromInput(input.StartAt, input.EndAt),
		},
		Type:  domain.MetricType(input.Type),
		Limit: int(input.Limit),
	})
	if err != nil {
		return nil, statsError(err, errorMessageMetricsLoadFailed)
	}

	return &websiteMetricsOutput{Body: newMetricResponses(metrics)}, nil
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

func newWebsiteStatsResponse(stats domain.WebsiteStats) WebsiteStatsResponse {
	return WebsiteStatsResponse{
		Pageviews:       stats.Pageviews,
		Visitors:        stats.Visitors,
		Visits:          stats.Visits,
		Bounces:         stats.Bounces,
		TotalTime:       stats.TotalTime,
		AvgVisitSeconds: stats.AvgVisitSeconds,
	}
}

func newPageviewResponses(buckets []domain.PageviewBucket) []PageviewResponse {
	result := make([]PageviewResponse, len(buckets))
	for i, bucket := range buckets {
		result[i] = PageviewResponse{
			Time:     bucket.Time,
			Label:    bucket.Label,
			Views:    bucket.Views,
			Visitors: bucket.Visitors,
		}
	}
	return result
}

func newMetricResponses(metrics []domain.Metric) []MetricResponse {
	result := make([]MetricResponse, len(metrics))
	for i, metric := range metrics {
		result[i] = MetricResponse{
			Name:     metric.Name,
			Views:    metric.Views,
			Visitors: metric.Visitors,
		}
	}
	return result
}
