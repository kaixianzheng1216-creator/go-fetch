package service

import (
	"context"
	"errors"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/repository"
)

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
	if err := service.requireWebsiteAccess(ctx, userID, websiteID); err != nil {
		return domain.WebsiteStats{}, err
	}

	start, end, _ := domain.DateRange(startAt, endAt, "")
	return service.store.WebsiteStats(ctx, websiteID, start, end)
}

func (service Stats) WebsitePageviews(ctx context.Context, userID, websiteID string, startAt, endAt *int64, unit string) ([]domain.PageviewPoint, error) {
	if err := service.requireWebsiteAccess(ctx, userID, websiteID); err != nil {
		return nil, err
	}

	start, end, parsedUnit := domain.DateRange(startAt, endAt, unit)
	return service.store.WebsitePageviews(ctx, websiteID, start, end, parsedUnit)
}

func (service Stats) WebsiteMetrics(ctx context.Context, userID, websiteID string, startAt, endAt *int64, metricType string, limit int) ([]domain.MetricRow, error) {
	if err := service.requireWebsiteAccess(ctx, userID, websiteID); err != nil {
		return nil, err
	}

	metric, isSupportedMetricType := domain.ParseMetricType(metricType)
	if !isSupportedMetricType {
		return nil, domain.ErrUnsupportedMetricType
	}

	if limit == 0 {
		limit = domain.DefaultMetricLimit
	}

	start, end, _ := domain.DateRange(startAt, endAt, "")
	return service.store.WebsiteMetrics(ctx, websiteID, start, end, metric, limit)
}

func (service Stats) requireWebsiteAccess(ctx context.Context, userID, websiteID string) error {
	if _, err := service.store.GetWebsite(ctx, userID, websiteID); err != nil {
		return WebsiteAccessError{Err: err}
	}
	return nil
}
