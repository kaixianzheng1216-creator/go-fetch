package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	storesqlc "github.com/kaixianzheng1216-creator/go-fetch/internal/repository/sqlc"
)

func (store *Store) WebsiteStats(ctx context.Context, websiteID uuid.UUID, start, end time.Time) (domain.WebsiteStats, error) {
	row, err := store.queries.WebsiteStats(ctx, storesqlc.WebsiteStatsParams{
		WebsiteID:         websiteID,
		StartAt:           start,
		EndAt:             end,
		PageviewEventType: int32(domain.EventTypePageView),
	})
	if err != nil {
		return domain.WebsiteStats{}, fmt.Errorf("load website stats: %w", err)
	}

	stats := domain.WebsiteStats{
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

func (store *Store) WebsitePageviews(ctx context.Context, websiteID uuid.UUID, start, end time.Time, unit domain.DateUnit) ([]domain.PageviewBucket, error) {
	rows, err := store.queries.Pageviews(ctx, storesqlc.PageviewsParams{
		Bucket:            domain.DateTruncUnit(unit),
		WebsiteID:         websiteID,
		StartAt:           start,
		EndAt:             end,
		PageviewEventType: int32(domain.EventTypePageView),
	})
	if err != nil {
		return nil, fmt.Errorf("load pageviews: %w", err)
	}

	buckets := make([]domain.PageviewBucket, 0, len(rows))
	for _, row := range rows {
		bucket := domain.PageviewBucket{
			Time:     row.Time,
			Views:    row.Views,
			Visitors: row.Visitors,
		}
		bucket.Label = domain.FormatBucket(bucket.Time, unit)
		buckets = append(buckets, bucket)
	}

	return buckets, nil
}

func (store *Store) WebsiteMetrics(ctx context.Context, websiteID uuid.UUID, start, end time.Time, metric domain.MetricType, limit int) ([]domain.Metric, error) {
	if _, isSupportedMetricType := domain.ParseMetricType(string(metric)); !isSupportedMetricType {
		return nil, domain.ErrUnsupportedMetricType
	}

	limit = domain.NormalizeMetricLimit(limit)
	if metric.IsSessionDimension() {
		rows, err := store.queries.SessionMetrics(ctx, storesqlc.SessionMetricsParams{
			Metric:     string(metric),
			WebsiteID:  websiteID,
			StartAt:    start,
			EndAt:      end,
			EventType:  int32(metric.EventType()),
			LimitCount: int32(limit),
		})
		if err != nil {
			return nil, fmt.Errorf("load session metrics: %w", err)
		}

		return toSessionMetrics(rows), nil
	}

	rows, err := store.queries.EventMetrics(ctx, storesqlc.EventMetricsParams{
		Metric:     string(metric),
		WebsiteID:  websiteID,
		StartAt:    start,
		EndAt:      end,
		EventType:  int32(metric.EventType()),
		LimitCount: int32(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("load event metrics: %w", err)
	}

	return toEventMetrics(rows), nil
}

func toSessionMetrics(rows []storesqlc.SessionMetricsRow) []domain.Metric {
	metrics := make([]domain.Metric, 0, len(rows))
	for _, row := range rows {
		metrics = append(metrics, domain.Metric{
			Name:     row.Name,
			Views:    row.Views,
			Visitors: row.Visitors,
		})
	}
	return metrics
}

func toEventMetrics(rows []storesqlc.EventMetricsRow) []domain.Metric {
	metrics := make([]domain.Metric, 0, len(rows))
	for _, row := range rows {
		metrics = append(metrics, domain.Metric{
			Name:     row.Name,
			Views:    row.Views,
			Visitors: row.Visitors,
		})
	}
	return metrics
}
