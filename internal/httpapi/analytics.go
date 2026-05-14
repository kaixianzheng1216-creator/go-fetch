package httpapi

import (
	"time"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
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

func WebsiteStatsFromDomain(stats domain.WebsiteStats) WebsiteStats {
	return WebsiteStats{
		Pageviews:       stats.Pageviews,
		Visitors:        stats.Visitors,
		Visits:          stats.Visits,
		Bounces:         stats.Bounces,
		TotalTime:       stats.TotalTime,
		AvgVisitSeconds: stats.AvgVisitSeconds,
	}
}

func PageviewPointFromDomain(point domain.PageviewPoint) PageviewPoint {
	return PageviewPoint{
		Time:     point.Time,
		Label:    point.Label,
		Views:    point.Views,
		Visitors: point.Visitors,
	}
}

func PageviewPointsFromDomain(points []domain.PageviewPoint) []PageviewPoint {
	result := make([]PageviewPoint, 0, len(points))
	for _, point := range points {
		result = append(result, PageviewPointFromDomain(point))
	}

	return result
}

func MetricRowFromDomain(row domain.MetricRow) MetricRow {
	return MetricRow{
		Name:     row.Name,
		Views:    row.Views,
		Visitors: row.Visitors,
	}
}

func MetricRowsFromDomain(rows []domain.MetricRow) []MetricRow {
	result := make([]MetricRow, 0, len(rows))
	for _, row := range rows {
		result = append(result, MetricRowFromDomain(row))
	}

	return result
}
