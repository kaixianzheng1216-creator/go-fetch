package repository

import (
	"context"
	"time"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
)

type AuthRepository interface {
	GetUserByUsername(ctx context.Context, username string) (domain.User, error)
}

type WebsiteRepository interface {
	ListWebsites(ctx context.Context, userID string) ([]domain.Website, error)
	CreateWebsite(ctx context.Context, userID, name, domainName string) (domain.Website, error)
	GetWebsite(ctx context.Context, userID, websiteID string) (domain.Website, error)
	UpdateWebsite(ctx context.Context, userID, websiteID, name, domainName string) error
	DeleteWebsite(ctx context.Context, userID, websiteID string) error
}

type TrackingRepository interface {
	GetWebsiteForCollection(ctx context.Context, websiteID string) (domain.Website, error)
	SaveEvent(ctx context.Context, event domain.EventInput) error
}

type StatsRepository interface {
	GetWebsite(ctx context.Context, userID, websiteID string) (domain.Website, error)
	WebsiteStats(ctx context.Context, websiteID string, start, end time.Time) (domain.WebsiteStats, error)
	WebsitePageviews(ctx context.Context, websiteID string, start, end time.Time, unit domain.DateUnit) ([]domain.PageviewPoint, error)
	WebsiteMetrics(ctx context.Context, websiteID string, start, end time.Time, metric domain.MetricType, limit int) ([]domain.MetricRow, error)
}
