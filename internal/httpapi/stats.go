package httpapi

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
)

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

type WebsiteStats struct {
	Pageviews       int64 `json:"pageviews"`
	Visitors        int64 `json:"visitors"`
	Visits          int64 `json:"visits"`
	Bounces         int64 `json:"bounces"`
	TotalTime       int64 `json:"totalTime"`
	AvgVisitSeconds int64 `json:"avgVisitSeconds"`
}

type PageviewPoint struct {
	Time     time.Time `json:"time"`
	Label    string    `json:"label"`
	Views    int64     `json:"views"`
	Visitors int64     `json:"visitors"`
}

type MetricRow struct {
	Name     string `json:"name"`
	Views    int64  `json:"views"`
	Visitors int64  `json:"visitors"`
}

type statsOutput struct {
	Body WebsiteStats
}

type pageviewsOutput struct {
	Body []PageviewPoint
}

type metricsOutput struct {
	Body []MetricRow
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

func (apiServer server) getWebsiteStats(ctx context.Context, input *statsRequest) (*statsOutput, error) {
	websiteID, err := parseUUID(input.WebsiteID, "websiteID")
	if err != nil {
		return nil, err
	}

	stats, err := apiServer.stats.WebsiteStats(ctx, currentUser(ctx).ID, websiteID, optionalTimeParam(input.StartAt), optionalTimeParam(input.EndAt))
	if err != nil {
		return nil, statsError(err, "加载统计数据失败")
	}

	return &statsOutput{Body: toWebsiteStatsResponse(stats)}, nil
}

func (apiServer server) getWebsitePageviews(ctx context.Context, input *pageviewsRequest) (*pageviewsOutput, error) {
	websiteID, err := parseUUID(input.WebsiteID, "websiteID")
	if err != nil {
		return nil, err
	}

	points, err := apiServer.stats.WebsitePageviews(ctx, currentUser(ctx).ID, websiteID, optionalTimeParam(input.StartAt), optionalTimeParam(input.EndAt), string(input.Unit))
	if err != nil {
		return nil, statsError(err, "加载页面浏览量失败")
	}

	return &pageviewsOutput{Body: toPageviewPointResponses(points)}, nil
}

func (apiServer server) getWebsiteMetrics(ctx context.Context, input *metricsRequest) (*metricsOutput, error) {
	websiteID, err := parseUUID(input.WebsiteID, "websiteID")
	if err != nil {
		return nil, err
	}

	rows, err := apiServer.stats.WebsiteMetrics(ctx, currentUser(ctx).ID, websiteID, optionalTimeParam(input.StartAt), optionalTimeParam(input.EndAt), string(input.Type), int(input.Limit))
	if err != nil {
		return nil, statsError(err, "加载指标数据失败")
	}

	return &metricsOutput{Body: toMetricRowResponses(rows)}, nil
}

func optionalTimeParam(value int64) *int64 {
	if value == 0 {
		return nil
	}
	return &value
}

func toWebsiteStatsResponse(stats domain.WebsiteStats) WebsiteStats {
	return WebsiteStats{
		Pageviews:       stats.Pageviews,
		Visitors:        stats.Visitors,
		Visits:          stats.Visits,
		Bounces:         stats.Bounces,
		TotalTime:       stats.TotalTime,
		AvgVisitSeconds: stats.AvgVisitSeconds,
	}
}

func toPageviewPointResponses(points []domain.PageviewPoint) []PageviewPoint {
	result := make([]PageviewPoint, 0, len(points))
	for _, point := range points {
		result = append(result, PageviewPoint{
			Time:     point.Time,
			Label:    point.Label,
			Views:    point.Views,
			Visitors: point.Visitors,
		})
	}
	return result
}

func toMetricRowResponses(rows []domain.MetricRow) []MetricRow {
	result := make([]MetricRow, 0, len(rows))
	for _, row := range rows {
		result = append(result, MetricRow{
			Name:     row.Name,
			Views:    row.Views,
			Visitors: row.Visitors,
		})
	}
	return result
}
