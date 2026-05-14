package store

import (
	"context"
	"fmt"
	"time"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	storedb "github.com/kaixianzheng1216-creator/go-fetch/internal/store/db"

	"github.com/google/uuid"
)

func (s *Store) WebsiteStats(ctx context.Context, websiteID string, start, end time.Time) (domain.WebsiteStats, error) {
	websiteUUID, err := uuid.Parse(websiteID)
	if err != nil {
		return domain.WebsiteStats{}, fmt.Errorf("解析网站 ID 失败: %w", err)
	}

	row, err := s.queries.WebsiteStats(ctx, storedb.WebsiteStatsParams{
		WebsiteID:         websiteUUID,
		StartAt:           start,
		EndAt:             end,
		PageviewEventType: int32(domain.EventTypePageView),
	})
	if err != nil {
		return domain.WebsiteStats{}, fmt.Errorf("加载网站统计失败: %w", err)
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

func (s *Store) Pageviews(ctx context.Context, websiteID string, start, end time.Time, unit domain.DateUnit) ([]domain.PageviewPoint, error) {
	websiteUUID, err := uuid.Parse(websiteID)
	if err != nil {
		return nil, fmt.Errorf("解析网站 ID 失败: %w", err)
	}

	rows, err := s.queries.Pageviews(ctx, storedb.PageviewsParams{
		Bucket:            domain.DateTruncUnit(unit),
		WebsiteID:         websiteUUID,
		StartAt:           start,
		EndAt:             end,
		PageviewEventType: int32(domain.EventTypePageView),
	})
	if err != nil {
		return nil, fmt.Errorf("加载浏览量失败: %w", err)
	}

	points := make([]domain.PageviewPoint, 0, len(rows))
	for _, row := range rows {
		point := domain.PageviewPoint{
			Time:     row.Time,
			Views:    row.Views,
			Visitors: row.Visitors,
		}
		point.Label = domain.FormatBucket(point.Time, unit)
		points = append(points, point)
	}

	return points, nil
}

func (s *Store) Metrics(ctx context.Context, websiteID string, start, end time.Time, metric domain.MetricType, limit int) ([]domain.MetricRow, error) {
	if _, ok := domain.ParseMetricType(string(metric)); !ok {
		return nil, domain.ErrUnsupportedMetricType
	}

	limit = domain.NormalizeMetricLimit(limit)
	websiteUUID, err := uuid.Parse(websiteID)
	if err != nil {
		return nil, fmt.Errorf("解析网站 ID 失败: %w", err)
	}

	if metric.IsSessionDimension() {
		rows, err := s.queries.SessionMetrics(ctx, storedb.SessionMetricsParams{
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

		metrics := make([]domain.MetricRow, 0, len(rows))
		for _, row := range rows {
			metrics = append(metrics, domain.MetricRow{
				Name:     row.Name,
				Views:    row.Views,
				Visitors: row.Visitors,
			})
		}

		return metrics, nil
	}

	rows, err := s.queries.EventMetrics(ctx, storedb.EventMetricsParams{
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

	metrics := make([]domain.MetricRow, 0, len(rows))
	for _, row := range rows {
		metrics = append(metrics, domain.MetricRow{
			Name:     row.Name,
			Views:    row.Views,
			Visitors: row.Visitors,
		})
	}

	return metrics, nil
}
