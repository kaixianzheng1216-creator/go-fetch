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
	WebsiteID uuid.UUID       `path:"websiteID" format:"uuid"`
	StartAt   int64           `query:"startAt"`
	EndAt     int64           `query:"endAt"`
	Type      MetricTypeParam `query:"type" required:"true"`
	Limit     MetricLimit     `query:"limit"`
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

type MetricLimit int

func (MetricLimit) Schema(huma.Registry) *huma.Schema {
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

func (apiServer server) registerStatsRoutes(humaAPI huma.API, authMiddleware huma.Middlewares) {
	huma.Register(
		humaAPI,
		securedOperation(http.MethodGet, "/api/websites/{websiteID}/stats", "websiteStats", "获取站点统计", "Analytics", authMiddleware),
		apiServer.getWebsiteStats,
	)

	huma.Register(
		humaAPI,
		securedOperation(http.MethodGet, "/api/websites/{websiteID}/pageviews", "websitePageviews", "获取页面浏览趋势", "Analytics", authMiddleware),
		apiServer.getWebsitePageviews,
	)

	huma.Register(
		humaAPI,
		securedOperation(http.MethodGet, "/api/websites/{websiteID}/metrics", "websiteMetrics", "获取站点指标", "Analytics", authMiddleware),
		apiServer.getWebsiteMetrics,
	)
}

func (apiServer server) getWebsiteStats(ctx context.Context, input *websiteStatsInput) (*websiteStatsOutput, error) {
	stats, err := apiServer.stats.Summary(ctx, service.StatsQuery{
		UserID:    currentUser(ctx).ID,
		WebsiteID: input.WebsiteID,
		Range:     dateRangeFromInput(input.StartAt, input.EndAt),
	})
	if err != nil {
		return nil, statsError(err, errorMessageStatsLoadFailed)
	}

	return &websiteStatsOutput{Body: toWebsiteStatsResponse(stats)}, nil
}

func (apiServer server) getWebsitePageviews(ctx context.Context, input *websitePageviewsInput) (*websitePageviewsOutput, error) {
	buckets, err := apiServer.stats.Pageviews(ctx, service.PageviewsQuery{
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

	return &websitePageviewsOutput{Body: toPageviewResponses(buckets)}, nil
}

func (apiServer server) getWebsiteMetrics(ctx context.Context, input *websiteMetricsInput) (*websiteMetricsOutput, error) {
	metricType, isSupportedMetricType := domain.ParseMetricType(string(input.Type))
	if !isSupportedMetricType {
		return nil, huma.Error400BadRequest(domain.ErrUnsupportedMetricType.Error())
	}

	metrics, err := apiServer.stats.Metrics(ctx, service.MetricsQuery{
		StatsQuery: service.StatsQuery{
			UserID:    currentUser(ctx).ID,
			WebsiteID: input.WebsiteID,
			Range:     dateRangeFromInput(input.StartAt, input.EndAt),
		},
		Type:  metricType,
		Limit: int(input.Limit),
	})
	if err != nil {
		return nil, statsError(err, errorMessageMetricsLoadFailed)
	}

	return &websiteMetricsOutput{Body: toMetricResponses(metrics)}, nil
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

func toWebsiteStatsResponse(stats domain.WebsiteStats) WebsiteStatsResponse {
	return WebsiteStatsResponse{
		Pageviews:       stats.Pageviews,
		Visitors:        stats.Visitors,
		Visits:          stats.Visits,
		Bounces:         stats.Bounces,
		TotalTime:       stats.TotalTime,
		AvgVisitSeconds: stats.AvgVisitSeconds,
	}
}

func toPageviewResponses(buckets []domain.PageviewBucket) []PageviewResponse {
	result := make([]PageviewResponse, 0, len(buckets))
	for _, bucket := range buckets {
		result = append(result, PageviewResponse{
			Time:     bucket.Time,
			Label:    bucket.Label,
			Views:    bucket.Views,
			Visitors: bucket.Visitors,
		})
	}
	return result
}

func toMetricResponses(metrics []domain.Metric) []MetricResponse {
	result := make([]MetricResponse, 0, len(metrics))
	for _, metric := range metrics {
		result = append(result, MetricResponse{
			Name:     metric.Name,
			Views:    metric.Views,
			Visitors: metric.Visitors,
		})
	}
	return result
}
