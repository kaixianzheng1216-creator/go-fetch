package store

import (
	"context"
	"fmt"
	"time"

	eventdomain "github.com/kaixianzheng1216-creator/go-fetch/internal/event"
	storesqlc "github.com/kaixianzheng1216-creator/go-fetch/internal/store/sqlc"

	"github.com/google/uuid"
)

func (s *Store) WebsiteStats(ctx context.Context, websiteID string, start, end time.Time) (eventdomain.WebsiteStats, error) {
	websiteUUID, err := uuid.Parse(websiteID)
	if err != nil {
		return eventdomain.WebsiteStats{}, fmt.Errorf("解析站点 ID 失败: %w", err)
	}

	row, err := s.queries.WebsiteStats(ctx, storesqlc.WebsiteStatsParams{
		WebsiteID:         websiteUUID,
		StartAt:           start,
		EndAt:             end,
		PageviewEventType: int32(eventdomain.EventTypePageView),
	})
	if err != nil {
		return eventdomain.WebsiteStats{}, fmt.Errorf("加载站点统计失败: %w", err)
	}

	stats := eventdomain.WebsiteStats{
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

func (s *Store) Pageviews(ctx context.Context, websiteID string, start, end time.Time, unit eventdomain.DateUnit) ([]eventdomain.PageviewPoint, error) {
	websiteUUID, err := uuid.Parse(websiteID)
	if err != nil {
		return nil, fmt.Errorf("解析站点 ID 失败: %w", err)
	}

	rows, err := s.queries.Pageviews(ctx, storesqlc.PageviewsParams{
		Bucket:            eventdomain.DateTruncUnit(unit),
		WebsiteID:         websiteUUID,
		StartAt:           start,
		EndAt:             end,
		PageviewEventType: int32(eventdomain.EventTypePageView),
	})
	if err != nil {
		return nil, fmt.Errorf("加载页面浏览量失败: %w", err)
	}

	points := make([]eventdomain.PageviewPoint, 0, len(rows))
	for _, row := range rows {
		point := eventdomain.PageviewPoint{
			Time:     row.Time,
			Views:    row.Views,
			Visitors: row.Visitors,
		}
		point.Label = eventdomain.FormatBucket(point.Time, unit)
		points = append(points, point)
	}

	return points, nil
}

func (s *Store) Metrics(ctx context.Context, websiteID string, start, end time.Time, metric eventdomain.MetricType, limit int) ([]eventdomain.MetricRow, error) {
	if _, ok := eventdomain.ParseMetricType(string(metric)); !ok {
		return nil, eventdomain.ErrUnsupportedMetricType
	}

	limit = eventdomain.NormalizeMetricLimit(limit)
	websiteUUID, err := uuid.Parse(websiteID)
	if err != nil {
		return nil, fmt.Errorf("解析站点 ID 失败: %w", err)
	}

	if metric.IsSessionDimension() {
		rows, err := s.queries.SessionMetrics(ctx, storesqlc.SessionMetricsParams{
			Metric:     string(metric),
			WebsiteID:  websiteUUID,
			StartAt:    start,
			EndAt:      end,
			EventType:  int32(metric.EventType()),
			LimitCount: int32(limit),
		})
		if err != nil {
			return nil, fmt.Errorf("加载会话指标失败: %w", err)
		}

		metrics := make([]eventdomain.MetricRow, 0, len(rows))
		for _, row := range rows {
			metrics = append(metrics, eventdomain.MetricRow{
				Name:     row.Name,
				Views:    row.Views,
				Visitors: row.Visitors,
			})
		}

		return metrics, nil
	}

	rows, err := s.queries.EventMetrics(ctx, storesqlc.EventMetricsParams{
		Metric:     string(metric),
		WebsiteID:  websiteUUID,
		StartAt:    start,
		EndAt:      end,
		EventType:  int32(metric.EventType()),
		LimitCount: int32(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("加载事件指标失败: %w", err)
	}

	metrics := make([]eventdomain.MetricRow, 0, len(rows))
	for _, row := range rows {
		metrics = append(metrics, eventdomain.MetricRow{
			Name:     row.Name,
			Views:    row.Views,
			Visitors: row.Visitors,
		})
	}

	return metrics, nil
}
