package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
)

var ErrInvalidDateRange = errors.New("startAt must be before or equal to endAt")

type StatsRepository interface {
	GetWebsite(ctx context.Context, userID, websiteID uuid.UUID) (domain.Website, error)
	WebsiteStats(ctx context.Context, websiteID uuid.UUID, start, end time.Time) (domain.WebsiteStats, error)
	WebsitePageviews(ctx context.Context, websiteID uuid.UUID, start, end time.Time, unit domain.DateUnit) ([]domain.PageviewBucket, error)
	WebsiteMetrics(ctx context.Context, websiteID uuid.UUID, start, end time.Time, metric domain.MetricType, limit int) ([]domain.Metric, error)
}

type DateRange struct {
	StartAt *time.Time
	EndAt   *time.Time
}

type StatsQuery struct {
	UserID    uuid.UUID
	WebsiteID uuid.UUID
	Range     DateRange
}

type PageviewsQuery struct {
	StatsQuery
	Unit domain.DateUnit
}

type MetricsQuery struct {
	StatsQuery
	Type  domain.MetricType
	Limit int
}

type StatsService struct {
	repository StatsRepository
	clock      clock
}

func NewStatsService(repository StatsRepository) StatsService {
	return StatsService{repository: repository, clock: systemClock}
}

func (svc StatsService) Summary(ctx context.Context, query StatsQuery) (domain.WebsiteStats, error) {
	start, end, err := svc.resolveStatsQuery(ctx, query)
	if err != nil {
		return domain.WebsiteStats{}, err
	}

	return svc.repository.WebsiteStats(ctx, query.WebsiteID, start, end)
}

func (svc StatsService) Pageviews(ctx context.Context, query PageviewsQuery) ([]domain.PageviewBucket, error) {
	start, end, err := svc.resolveStatsQuery(ctx, query.StatsQuery)
	if err != nil {
		return nil, err
	}

	unit, isSupportedDateUnit := domain.ParseDateUnit(string(query.Unit))
	if !isSupportedDateUnit {
		return nil, domain.ErrUnsupportedDateUnit
	}

	return svc.repository.WebsitePageviews(ctx, query.WebsiteID, start, end, unit)
}

func (svc StatsService) Metrics(ctx context.Context, query MetricsQuery) ([]domain.Metric, error) {
	start, end, err := svc.resolveStatsQuery(ctx, query.StatsQuery)
	if err != nil {
		return nil, err
	}

	if _, isSupportedMetricType := domain.ParseMetricType(string(query.Type)); !isSupportedMetricType {
		return nil, domain.ErrUnsupportedMetricType
	}

	return svc.repository.WebsiteMetrics(ctx, query.WebsiteID, start, end, query.Type, domain.NormalizeMetricLimit(query.Limit))
}

func (svc StatsService) resolveStatsQuery(ctx context.Context, query StatsQuery) (time.Time, time.Time, error) {
	start, end, err := statsDateRange(svc.now(), query.Range)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	if err := svc.requireWebsiteAccess(ctx, query.UserID, query.WebsiteID); err != nil {
		return time.Time{}, time.Time{}, err
	}

	return start, end, nil
}

func (svc StatsService) requireWebsiteAccess(ctx context.Context, userID, websiteID uuid.UUID) error {
	if _, err := svc.repository.GetWebsite(ctx, userID, websiteID); err != nil {
		return websiteAccessError{err: err}
	}
	return nil
}

func (svc StatsService) now() time.Time {
	if svc.clock == nil {
		return systemClock()
	}
	return svc.clock()
}

func IsWebsiteAccessError(err error) bool {
	var accessError websiteAccessError
	return errors.As(err, &accessError)
}

type websiteAccessError struct {
	err error
}

func (err websiteAccessError) Error() string {
	return "website access: " + err.err.Error()
}

func (err websiteAccessError) Unwrap() error {
	return err.err
}

func statsDateRange(now time.Time, dateRange DateRange) (time.Time, time.Time, error) {
	start := now.Add(-domain.DefaultDateLookback)
	end := now
	if dateRange.StartAt != nil {
		start = *dateRange.StartAt
	}
	if dateRange.EndAt != nil {
		end = *dateRange.EndAt
	}

	if start.After(end) {
		return time.Time{}, time.Time{}, ErrInvalidDateRange
	}

	return start, end, nil
}
