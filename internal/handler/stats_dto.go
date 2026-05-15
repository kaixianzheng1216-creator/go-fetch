package handler

import (
	"time"

	"github.com/danielgtaylor/huma/v2"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/model"
)

type DateUnitParam string

func (DateUnitParam) Schema(huma.Registry) *huma.Schema {
	return &huma.Schema{
		Type:    huma.TypeString,
		Enum:    enumValues(model.DateUnitValues()),
		Default: string(model.DefaultDateUnit),
	}
}

type MetricTypeParam string

func (MetricTypeParam) Schema(huma.Registry) *huma.Schema {
	return &huma.Schema{
		Type: huma.TypeString,
		Enum: enumValues(model.MetricTypeValues()),
	}
}

type MetricLimit int

func (MetricLimit) Schema(huma.Registry) *huma.Schema {
	minValue := 1.0
	maxValue := float64(model.MaxMetricLimit)
	return &huma.Schema{
		Type:    huma.TypeInteger,
		Default: model.DefaultMetricLimit,
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

type StatsOutput struct {
	Body WebsiteStats
}

type PageviewsOutput struct {
	Body []PageviewPoint
}

type MetricsOutput struct {
	Body []MetricRow
}

func NewStatsOutput(stats WebsiteStats) *StatsOutput {
	return &StatsOutput{Body: stats}
}

func NewPageviewsOutput(points []PageviewPoint) *PageviewsOutput {
	return &PageviewsOutput{Body: points}
}

func NewMetricsOutput(rows []MetricRow) *MetricsOutput {
	return &MetricsOutput{Body: rows}
}

func OptionalTimeParam(value int64) *int64 {
	if value == 0 {
		return nil
	}
	return &value
}

func ToWebsiteStats(stats model.WebsiteStats) WebsiteStats {
	return WebsiteStats{
		Pageviews:       stats.Pageviews,
		Visitors:        stats.Visitors,
		Visits:          stats.Visits,
		Bounces:         stats.Bounces,
		TotalTime:       stats.TotalTime,
		AvgVisitSeconds: stats.AvgVisitSeconds,
	}
}

func ToPageviewPoint(point model.PageviewPoint) PageviewPoint {
	return PageviewPoint{
		Time:     point.Time,
		Label:    point.Label,
		Views:    point.Views,
		Visitors: point.Visitors,
	}
}

func ToPageviewPoints(points []model.PageviewPoint) []PageviewPoint {
	result := make([]PageviewPoint, 0, len(points))
	for _, point := range points {
		result = append(result, ToPageviewPoint(point))
	}

	return result
}

func ToMetricRow(row model.MetricRow) MetricRow {
	return MetricRow{
		Name:     row.Name,
		Views:    row.Views,
		Visitors: row.Visitors,
	}
}

func ToMetricRows(rows []model.MetricRow) []MetricRow {
	result := make([]MetricRow, 0, len(rows))
	for _, row := range rows {
		result = append(result, ToMetricRow(row))
	}

	return result
}
