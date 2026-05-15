package service

import (
	"context"
	"errors"
	"time"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/repository"
)

var ErrInvalidDateRange = errors.New("startAt must be before or equal to endAt")

type Stats struct {
	store repository.StatsRepository
}

type WebsiteAccessError struct {
	Err error
}

func (err WebsiteAccessError) Error() string {
	return err.Err.Error()
}

func (err WebsiteAccessError) Unwrap() error {
	return err.Err
}

func IsWebsiteAccessError(err error) bool {
	var accessError WebsiteAccessError
	return errors.As(err, &accessError)
}

func NewStats(store repository.StatsRepository) Stats {
	return Stats{store: store}
}

func (service Stats) WebsiteStats(ctx context.Context, userID, websiteID string, startAt, endAt *int64) (domain.WebsiteStats, error) {
	start, end, _, err := statsDateRange(startAt, endAt, "")
	if err != nil {
		return domain.WebsiteStats{}, err
	}

	if err := service.requireWebsiteAccess(ctx, userID, websiteID); err != nil {
		return domain.WebsiteStats{}, err
	}

	return service.store.WebsiteStats(ctx, websiteID, start, end)
}

func (service Stats) WebsitePageviews(ctx context.Context, userID, websiteID string, startAt, endAt *int64, unit string) ([]domain.PageviewPoint, error) {
	start, end, parsedUnit, err := statsDateRange(startAt, endAt, unit)
	if err != nil {
		return nil, err
	}

	if err := service.requireWebsiteAccess(ctx, userID, websiteID); err != nil {
		return nil, err
	}

	return service.store.WebsitePageviews(ctx, websiteID, start, end, parsedUnit)
}

func (service Stats) WebsiteMetrics(ctx context.Context, userID, websiteID string, startAt, endAt *int64, metricType string, limit int) ([]domain.MetricRow, error) {
	start, end, _, err := statsDateRange(startAt, endAt, "")
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

func (service Stats) requireWebsiteAccess(ctx context.Context, userID, websiteID string) error {
	if _, err := service.store.GetWebsite(ctx, userID, websiteID); err != nil {
		return WebsiteAccessError{Err: err}
	}
	return nil
}

func statsDateRange(startAt, endAt *int64, unit string) (time.Time, time.Time, domain.DateUnit, error) {
	start, end, parsedUnit := domain.DateRange(startAt, endAt, unit)
	if start.After(end) {
		return time.Time{}, time.Time{}, "", ErrInvalidDateRange
	}

	return start, end, parsedUnit, nil
}
