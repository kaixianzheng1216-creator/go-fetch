package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
)

var ErrInvalidDateRange = errors.New("startAt must be before or equal to endAt")

type StatsStore interface {
	GetWebsite(ctx context.Context, userID, websiteID uuid.UUID) (domain.Website, error)
	WebsiteStats(ctx context.Context, websiteID uuid.UUID, start, end time.Time) (domain.WebsiteStats, error)
	WebsitePageviews(ctx context.Context, websiteID uuid.UUID, start, end time.Time, unit domain.DateUnit) ([]domain.PageviewPoint, error)
	WebsiteMetrics(ctx context.Context, websiteID uuid.UUID, start, end time.Time, metric domain.MetricType, limit int) ([]domain.MetricRow, error)
}

type websiteAccessError struct {
	err error
}

func (err websiteAccessError) Error() string {
	return err.err.Error()
}

func (err websiteAccessError) Unwrap() error {
	return err.err
}

func IsWebsiteAccessError(err error) bool {
	var accessError websiteAccessError
	return errors.As(err, &accessError)
}

type Stats struct {
	store StatsStore
	clock Clock
}

func NewStats(store StatsStore) Stats {
	return Stats{store: store, clock: systemClock}
}

func (service Stats) WebsiteStats(ctx context.Context, userID, websiteID uuid.UUID, startAt, endAt *int64) (domain.WebsiteStats, error) {
	start, end, _, err := statsDateRange(service.now(), startAt, endAt, "")
	if err != nil {
		return domain.WebsiteStats{}, err
	}

	if err := service.requireWebsiteAccess(ctx, userID, websiteID); err != nil {
		return domain.WebsiteStats{}, err
	}

	return service.store.WebsiteStats(ctx, websiteID, start, end)
}

func (service Stats) WebsitePageviews(ctx context.Context, userID, websiteID uuid.UUID, startAt, endAt *int64, unit string) ([]domain.PageviewPoint, error) {
	start, end, parsedUnit, err := statsDateRange(service.now(), startAt, endAt, unit)
	if err != nil {
		return nil, err
	}

	if err := service.requireWebsiteAccess(ctx, userID, websiteID); err != nil {
		return nil, err
	}

	return service.store.WebsitePageviews(ctx, websiteID, start, end, parsedUnit)
}

func (service Stats) WebsiteMetrics(ctx context.Context, userID, websiteID uuid.UUID, startAt, endAt *int64, metricType string, limit int) ([]domain.MetricRow, error) {
	start, end, _, err := statsDateRange(service.now(), startAt, endAt, "")
	if err != nil {
		return nil, err
	}

	if err := service.requireWebsiteAccess(ctx, userID, websiteID); err != nil {
		return nil, err
	}

	metric, isSupportedMetricType := domain.ParseMetricType(metricType)
	if !isSupportedMetricType {
		return nil, domain.ErrUnsupportedMetricType
	}

	return service.store.WebsiteMetrics(ctx, websiteID, start, end, metric, domain.NormalizeMetricLimit(limit))
}

func (service Stats) requireWebsiteAccess(ctx context.Context, userID, websiteID uuid.UUID) error {
	if _, err := service.store.GetWebsite(ctx, userID, websiteID); err != nil {
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

func statsDateRange(now time.Time, startAt, endAt *int64, unit string) (time.Time, time.Time, domain.DateUnit, error) {
	start := now.Add(-domain.DefaultDateLookback)
	end := now
	if startAt != nil {
		start = time.UnixMilli(*startAt)
	}
	if endAt != nil {
		end = time.UnixMilli(*endAt)
	}

	parsedUnit := domain.ParseDateUnit(unit)
	if start.After(end) {
		return time.Time{}, time.Time{}, "", ErrInvalidDateRange
	}

	return start, end, parsedUnit, nil
}
