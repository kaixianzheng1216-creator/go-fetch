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

type Websites struct {
	store WebsiteStore
}

func NewWebsites(store WebsiteStore) Websites {
	return Websites{store: store}
}

func (service Websites) List(ctx context.Context, userID uuid.UUID) ([]domain.Website, error) {
	return service.store.ListWebsites(ctx, userID)
}

func (service Websites) Create(ctx context.Context, userID uuid.UUID, name, domainName string) (domain.Website, error) {
	name, domainName = normalizeWebsiteInput(name, domainName)
	if name == "" {
		return domain.Website{}, ErrInvalidWebsiteName
	}
	return service.store.CreateWebsite(ctx, userID, name, domainName)
}

func (service Websites) Get(ctx context.Context, userID, websiteID uuid.UUID) (domain.Website, error) {
	return service.store.GetWebsite(ctx, userID, websiteID)
}

func (service Websites) Update(ctx context.Context, userID, websiteID uuid.UUID, name, domainName string) (domain.Website, error) {
	name, domainName = normalizeWebsiteInput(name, domainName)
	if name == "" {
		return domain.Website{}, ErrInvalidWebsiteName
	}

	if err := service.store.UpdateWebsite(ctx, userID, websiteID, name, domainName); err != nil {
		return domain.Website{}, err
	}
	return service.store.GetWebsite(ctx, userID, websiteID)
}

func (service Websites) Delete(ctx context.Context, userID, websiteID uuid.UUID) error {
	return service.store.DeleteWebsite(ctx, userID, websiteID)
}

func normalizeWebsiteInput(name, domain string) (string, string) {
	return strings.TrimSpace(name), strings.TrimSpace(domain)
}
