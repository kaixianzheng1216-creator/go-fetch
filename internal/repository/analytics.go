package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/model"
	storesqlc "github.com/kaixianzheng1216-creator/go-fetch/internal/repository/sqlc"
)

func (store *Store) WebsiteStats(ctx context.Context, websiteID string, start, end time.Time) (model.WebsiteStats, error) {
	websiteUUID, err := uuid.Parse(websiteID)
	if err != nil {
		return model.WebsiteStats{}, fmt.Errorf("parse website ID: %w", err)
	}

	row, err := store.queries.WebsiteStats(ctx, storesqlc.WebsiteStatsParams{
		WebsiteID:         websiteUUID,
		StartAt:           start,
		EndAt:             end,
		PageviewEventType: int32(model.EventTypePageView),
	})
	if err != nil {
		return model.WebsiteStats{}, fmt.Errorf("load website stats: %w", err)
	}

	stats := model.WebsiteStats{
		Pageviews: row.Pageviews,
		Visitors:  row.Visitors,
		Visits:    row.Visits,
		Bounces:   row.Bounces,
		TotalTime: row.TotalTime,
	}
	if stats.Visits > 0 {
		stats.AvgVisitSeconds = stats.TotalTime / stats.Visits
	}

	return stats, nil
}

func (store *Store) WebsitePageviews(ctx context.Context, websiteID string, start, end time.Time, unit model.DateUnit) ([]model.PageviewPoint, error) {
	websiteUUID, err := uuid.Parse(websiteID)
	if err != nil {
		return nil, fmt.Errorf("parse website ID: %w", err)
	}

	rows, err := store.queries.Pageviews(ctx, storesqlc.PageviewsParams{
		Bucket:            model.DateTruncUnit(unit),
		WebsiteID:         websiteUUID,
		StartAt:           start,
		EndAt:             end,
		PageviewEventType: int32(model.EventTypePageView),
	})
	if err != nil {
		return nil, fmt.Errorf("load pageviews: %w", err)
	}

	points := make([]model.PageviewPoint, 0, len(rows))
	for _, row := range rows {
		point := model.PageviewPoint{
			Time:     row.Time,
			Views:    row.Views,
			Visitors: row.Visitors,
		}
		point.Label = model.FormatBucket(point.Time, unit)
		points = append(points, point)
	}

	return points, nil
}

func (store *Store) WebsiteMetrics(ctx context.Context, websiteID string, start, end time.Time, metric model.MetricType, limit int) ([]model.MetricRow, error) {
	if _, isSupportedMetricType := model.ParseMetricType(string(metric)); !isSupportedMetricType {
		return nil, model.ErrUnsupportedMetricType
	}

	limit = model.NormalizeMetricLimit(limit)
	websiteUUID, err := uuid.Parse(websiteID)
	if err != nil {
		return nil, fmt.Errorf("parse website ID: %w", err)
	}

	if metric.IsSessionDimension() {
		rows, err := store.queries.SessionMetrics(ctx, storesqlc.SessionMetricsParams{
			Metric:     string(metric),
			WebsiteID:  websiteUUID,
			StartAt:    start,
			EndAt:      end,
			EventType:  int32(metric.EventType()),
			LimitCount: int32(limit),
		})
		if err != nil {
			return nil, fmt.Errorf("load session metrics: %w", err)
		}

		metrics := make([]model.MetricRow, 0, len(rows))
		for _, row := range rows {
			metrics = append(metrics, model.MetricRow{
				Name:     row.Name,
				Views:    row.Views,
				Visitors: row.Visitors,
			})
		}
		return metrics, nil
	}

	rows, err := store.queries.EventMetrics(ctx, storesqlc.EventMetricsParams{
		Metric:     string(metric),
		WebsiteID:  websiteUUID,
		StartAt:    start,
		EndAt:      end,
		EventType:  int32(metric.EventType()),
		LimitCount: int32(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("load event metrics: %w", err)
	}

	metrics := make([]model.MetricRow, 0, len(rows))
	for _, row := range rows {
		metrics = append(metrics, model.MetricRow{
			Name:     row.Name,
			Views:    row.Views,
			Visitors: row.Visitors,
		})
	}

	return metrics, nil
}
