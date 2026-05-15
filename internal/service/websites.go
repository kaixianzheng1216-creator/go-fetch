package service

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
)

var ErrInvalidWebsiteName = errors.New("website name cannot be empty")

type WebsiteStore interface {
	ListWebsites(ctx context.Context, userID uuid.UUID) ([]domain.Website, error)
	CreateWebsite(ctx context.Context, userID uuid.UUID, name, domainName string) (domain.Website, error)
	GetWebsite(ctx context.Context, userID, websiteID uuid.UUID) (domain.Website, error)
	UpdateWebsite(ctx context.Context, userID, websiteID uuid.UUID, name, domainName string) error
	DeleteWebsite(ctx context.Context, userID, websiteID uuid.UUID) error
}

type Website struct {
	store WebsiteStore
}

func NewWebsite(store WebsiteStore) Website {
	return Website{store: store}
}

func (service Website) List(ctx context.Context, userID uuid.UUID) ([]domain.Website, error) {
	return service.store.ListWebsites(ctx, userID)
}

func (service Website) Create(ctx context.Context, userID uuid.UUID, name, domainName string) (domain.Website, error) {
	name, domainName = normalizeWebsiteInput(name, domainName)
	if name == "" {
		return domain.Website{}, ErrInvalidWebsiteName
	}
	return service.store.CreateWebsite(ctx, userID, name, domainName)
}

func (service Website) Get(ctx context.Context, userID, websiteID uuid.UUID) (domain.Website, error) {
	return service.store.GetWebsite(ctx, userID, websiteID)
}

func (service Website) Update(ctx context.Context, userID, websiteID uuid.UUID, name, domainName string) (domain.Website, error) {
	name, domainName = normalizeWebsiteInput(name, domainName)
	if name == "" {
		return domain.Website{}, ErrInvalidWebsiteName
	}

	if err := service.store.UpdateWebsite(ctx, userID, websiteID, name, domainName); err != nil {
		return domain.Website{}, err
	}
	return service.store.GetWebsite(ctx, userID, websiteID)
}

func (service Website) Delete(ctx context.Context, userID, websiteID uuid.UUID) error {
	return service.store.DeleteWebsite(ctx, userID, websiteID)
}

func normalizeWebsiteInput(name, domain string) (string, string) {
	return strings.TrimSpace(name), strings.TrimSpace(domain)
}
