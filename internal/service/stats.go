package service

import (
	"context"
	"errors"
	"time"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/model"
)

type StatsStore interface {
	GetWebsite(ctx context.Context, userID, websiteID string) (model.Website, error)
	WebsiteStats(ctx context.Context, websiteID string, start, end time.Time) (model.WebsiteStats, error)
	WebsitePageviews(ctx context.Context, websiteID string, start, end time.Time, unit model.DateUnit) ([]model.PageviewPoint, error)
	WebsiteMetrics(ctx context.Context, websiteID string, start, end time.Time, metric model.MetricType, limit int) ([]model.MetricRow, error)
}

type Stats struct {
	store StatsStore
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

func NewStats(store StatsStore) Stats {
	return Stats{store: store}
}

func (service Stats) WebsiteStats(ctx context.Context, userID, websiteID string, startAt, endAt *int64) (model.WebsiteStats, error) {
	if err := service.requireWebsiteAccess(ctx, userID, websiteID); err != nil {
		return model.WebsiteStats{}, err
	}

	start, end, _ := model.DateRange(startAt, endAt, "")
	return service.store.WebsiteStats(ctx, websiteID, start, end)
}

func (service Stats) WebsitePageviews(ctx context.Context, userID, websiteID string, startAt, endAt *int64, unit string) ([]model.PageviewPoint, error) {
	if err := service.requireWebsiteAccess(ctx, userID, websiteID); err != nil {
		return nil, err
	}

	start, end, parsedUnit := model.DateRange(startAt, endAt, unit)
	return service.store.WebsitePageviews(ctx, websiteID, start, end, parsedUnit)
}

func (service Stats) WebsiteMetrics(ctx context.Context, userID, websiteID string, startAt, endAt *int64, metricType string, limit int) ([]model.MetricRow, error) {
	if err := service.requireWebsiteAccess(ctx, userID, websiteID); err != nil {
		return nil, err
	}

	metric, isSupportedMetricType := model.ParseMetricType(metricType)
	if !isSupportedMetricType {
		return nil, model.ErrUnsupportedMetricType
	}

	if limit == 0 {
		limit = model.DefaultMetricLimit
	}

	start, end, _ := model.DateRange(startAt, endAt, "")
	return service.store.WebsiteMetrics(ctx, websiteID, start, end, metric, limit)
}

func (service Stats) requireWebsiteAccess(ctx context.Context, userID, websiteID string) error {
	if _, err := service.store.GetWebsite(ctx, userID, websiteID); err != nil {
		return WebsiteAccessError{Err: err}
	}
	return nil
}
