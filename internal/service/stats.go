package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
)

var ErrInvalidDateRange = errors.New("startAt must be before or equal to endAt")

// AnalyticsRepository reads analytics data for the stats service.
type AnalyticsRepository interface {
	GetWebsite(ctx context.Context, userID, websiteID uuid.UUID) (domain.Website, error)
	WebsiteStats(ctx context.Context, websiteID uuid.UUID, start, end time.Time) (domain.WebsiteStats, error)
	WebsitePageviews(ctx context.Context, websiteID uuid.UUID, start, end time.Time, unit domain.DateUnit) ([]domain.PageviewBucket, error)
	WebsiteMetrics(ctx context.Context, websiteID uuid.UUID, start, end time.Time, metric domain.MetricType, limit int) ([]domain.Metric, error)
}

// DateRange optionally constrains analytics queries by time.
type DateRange struct {
	StartAt *time.Time
	EndAt   *time.Time
}

// StatsParams scopes a website analytics request to a user-owned website.
type StatsParams struct {
	UserID    uuid.UUID
	WebsiteID uuid.UUID
	Range     DateRange
}

// PageviewsParams requests pageview buckets for a website.
type PageviewsParams struct {
	StatsParams
	Unit domain.DateUnit
}

// MetricsParams requests top metrics for a website.
type MetricsParams struct {
	StatsParams
	Type  domain.MetricType
	Limit int
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

func IsWebsiteAccessError(err error) bool {
	var accessError websiteAccessError
	return errors.As(err, &accessError)
}

type Stats struct {
	repository AnalyticsRepository
	clock      Clock
}

func NewStats(repository AnalyticsRepository) Stats {
	return Stats{repository: repository, clock: systemClock}
}

func (service Stats) Summary(ctx context.Context, params StatsParams) (domain.WebsiteStats, error) {
	start, end, err := statsDateRange(service.now(), params.Range)
	if err != nil {
		return domain.WebsiteStats{}, err
	}

	if err := service.requireWebsiteAccess(ctx, params.UserID, params.WebsiteID); err != nil {
		return domain.WebsiteStats{}, err
	}

	return service.repository.WebsiteStats(ctx, params.WebsiteID, start, end)
}

func (service Stats) Pageviews(ctx context.Context, params PageviewsParams) ([]domain.PageviewBucket, error) {
	start, end, err := statsDateRange(service.now(), params.Range)
	if err != nil {
		return nil, err
	}

	if err := service.requireWebsiteAccess(ctx, params.UserID, params.WebsiteID); err != nil {
		return nil, err
	}

	return service.repository.WebsitePageviews(ctx, params.WebsiteID, start, end, domain.ParseDateUnit(string(params.Unit)))
}

func (service Stats) Metrics(ctx context.Context, params MetricsParams) ([]domain.Metric, error) {
	start, end, err := statsDateRange(service.now(), params.Range)
	if err != nil {
		return nil, err
	}

	if err := service.requireWebsiteAccess(ctx, params.UserID, params.WebsiteID); err != nil {
		return nil, err
	}

	if _, isSupportedMetricType := domain.ParseMetricType(string(params.Type)); !isSupportedMetricType {
		return nil, domain.ErrUnsupportedMetricType
	}

	return service.repository.WebsiteMetrics(ctx, params.WebsiteID, start, end, params.Type, domain.NormalizeMetricLimit(params.Limit))
}

func (service Stats) requireWebsiteAccess(ctx context.Context, userID, websiteID uuid.UUID) error {
	if _, err := service.repository.GetWebsite(ctx, userID, websiteID); err != nil {
		return websiteAccessError{err: err}
	}
	return nil
}

func (service Stats) now() time.Time {
	if service.clock == nil {
		return systemClock()
	}
	return service.clock()
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
