package httpapi

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/service"
)

type statsInput struct {
	WebsiteID string `path:"websiteID" format:"uuid"`
	StartAt   int64  `query:"startAt"`
	EndAt     int64  `query:"endAt"`
}

type pageviewsInput struct {
	WebsiteID string        `path:"websiteID" format:"uuid"`
	StartAt   int64         `query:"startAt"`
	EndAt     int64         `query:"endAt"`
	Unit      DateUnitParam `query:"unit"`
}

type metricsInput struct {
	WebsiteID string          `path:"websiteID" format:"uuid"`
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

type PageviewPointResponse struct {
	Time     time.Time `json:"time"`
	Label    string    `json:"label"`
	Views    int64     `json:"views"`
	Visitors int64     `json:"visitors"`
}

type MetricRowResponse struct {
	Name     string `json:"name"`
	Views    int64  `json:"views"`
	Visitors int64  `json:"visitors"`
}

type statsOutput struct {
	Body WebsiteStatsResponse
}

type pageviewsOutput struct {
	Body []PageviewPointResponse
}

type metricsOutput struct {
	Body []MetricRowResponse
}

func (apiServer server) registerStatsRoutes(humaAPI huma.API, authMiddleware huma.Middlewares) {
	huma.Register(
		humaAPI,
		huma.Operation{
			Method:      http.MethodGet,
			Path:        "/api/websites/{websiteID}/stats",
			OperationID: "websiteStats",
			Summary:     "获取站点统计",
			Tags:        []string{"Analytics"},
			Security:    []map[string][]string{{"sessionCookie": {}}},
			Middlewares: authMiddleware,
		},
		apiServer.getWebsiteStats,
	)

	huma.Register(
		humaAPI,
		huma.Operation{
			Method:      http.MethodGet,
			Path:        "/api/websites/{websiteID}/pageviews",
			OperationID: "websitePageviews",
			Summary:     "获取页面浏览趋势",
			Tags:        []string{"Analytics"},
			Security:    []map[string][]string{{"sessionCookie": {}}},
			Middlewares: authMiddleware,
		},
		apiServer.getWebsitePageviews,
	)

	huma.Register(
		humaAPI,
		huma.Operation{
			Method:      http.MethodGet,
			Path:        "/api/websites/{websiteID}/metrics",
			OperationID: "websiteMetrics",
			Summary:     "获取站点指标",
			Tags:        []string{"Analytics"},
			Security:    []map[string][]string{{"sessionCookie": {}}},
			Middlewares: authMiddleware,
		},
		apiServer.getWebsiteMetrics,
	)
}

func (apiServer server) getWebsiteStats(ctx context.Context, input *statsInput) (*statsOutput, error) {
	websiteID, err := parseUUID(input.WebsiteID, "websiteID")
	if err != nil {
		return nil, err
	}

	stats, err := apiServer.stats.Summary(ctx, service.StatsQuery{
		UserID:    currentUser(ctx).ID,
		WebsiteID: websiteID,
		Range:     dateRangeFromInput(input.StartAt, input.EndAt),
	})
	if err != nil {
		return nil, statsError(err, "加载统计数据失败")
	}

	return &statsOutput{Body: toWebsiteStatsResponse(stats)}, nil
}

func (apiServer server) getWebsitePageviews(ctx context.Context, input *pageviewsInput) (*pageviewsOutput, error) {
	websiteID, err := parseUUID(input.WebsiteID, "websiteID")
	if err != nil {
		return nil, err
	}

	points, err := apiServer.stats.Pageviews(ctx, service.PageviewsQuery{
		StatsQuery: service.StatsQuery{
			UserID:    currentUser(ctx).ID,
			WebsiteID: websiteID,
			Range:     dateRangeFromInput(input.StartAt, input.EndAt),
		},
		Unit: domain.ParseDateUnit(string(input.Unit)),
	})
	if err != nil {
		return nil, statsError(err, "加载页面浏览量失败")
	}

	return &pageviewsOutput{Body: toPageviewPointResponses(points)}, nil
}

func (apiServer server) getWebsiteMetrics(ctx context.Context, input *metricsInput) (*metricsOutput, error) {
	websiteID, err := parseUUID(input.WebsiteID, "websiteID")
	if err != nil {
		return nil, err
	}

	metricType, isSupportedMetricType := domain.ParseMetricType(string(input.Type))
	if !isSupportedMetricType {
		return nil, huma.Error400BadRequest(domain.ErrUnsupportedMetricType.Error())
	}

	rows, err := apiServer.stats.Metrics(ctx, service.MetricsQuery{
		StatsQuery: service.StatsQuery{
			UserID:    currentUser(ctx).ID,
			WebsiteID: websiteID,
			Range:     dateRangeFromInput(input.StartAt, input.EndAt),
		},
		Type:  metricType,
		Limit: int(input.Limit),
	})
	if err != nil {
		return nil, statsError(err, "加载指标数据失败")
	}

	return &metricsOutput{Body: toMetricRowResponses(rows)}, nil
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

func toPageviewPointResponses(points []domain.PageviewPoint) []PageviewPointResponse {
	result := make([]PageviewPointResponse, 0, len(points))
	for _, point := range points {
		result = append(result, PageviewPointResponse{
			Time:     point.Time,
			Label:    point.Label,
			Views:    point.Views,
			Visitors: point.Visitors,
		})
	}
	return result
}

func toMetricRowResponses(rows []domain.MetricRow) []MetricRowResponse {
	result := make([]MetricRowResponse, 0, len(rows))
	for _, row := range rows {
		result = append(result, MetricRowResponse{
			Name:     row.Name,
			Views:    row.Views,
			Visitors: row.Visitors,
		})
	}
	return result
}
