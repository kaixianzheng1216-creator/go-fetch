package summary

import (
	"time"

	eventdomain "github.com/kaixianzheng1216-creator/go-fetch/internal/domain/event"
)

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

func newStatsOutput(stats WebsiteStats) *statsOutput {
	return &statsOutput{Body: stats}
}

func newPageviewsOutput(points []PageviewPoint) *pageviewsOutput {
	return &pageviewsOutput{Body: points}
}

func newMetricsOutput(rows []MetricRow) *metricsOutput {
	return &metricsOutput{Body: rows}
}

func ToWebsiteStats(stats eventdomain.WebsiteStats) WebsiteStats {
	return WebsiteStats{
		Pageviews:       stats.Pageviews,
		Visitors:        stats.Visitors,
		Visits:          stats.Visits,
		Bounces:         stats.Bounces,
		TotalTime:       stats.TotalTime,
		AvgVisitSeconds: stats.AvgVisitSeconds,
	}
}

func ToPageviewPoint(point eventdomain.PageviewPoint) PageviewPoint {
	return PageviewPoint{
		Time:     point.Time,
		Label:    point.Label,
		Views:    point.Views,
		Visitors: point.Visitors,
	}
}

func ToPageviewPoints(points []eventdomain.PageviewPoint) []PageviewPoint {
	result := make([]PageviewPoint, 0, len(points))
	for _, point := range points {
		result = append(result, ToPageviewPoint(point))
	}

	return result
}

func ToMetricRow(row eventdomain.MetricRow) MetricRow {
	return MetricRow{
		Name:     row.Name,
		Views:    row.Views,
		Visitors: row.Visitors,
	}
}

func ToMetricRows(rows []eventdomain.MetricRow) []MetricRow {
	result := make([]MetricRow, 0, len(rows))
	for _, row := range rows {
		result = append(result, ToMetricRow(row))
	}

	return result
}
