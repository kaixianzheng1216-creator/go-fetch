package service

import (
	"context"
	"errors"
	"strings"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/repository"
)

var ErrInvalidWebsiteName = errors.New("website name cannot be empty")

type Website struct {
	store repository.WebsiteRepository
}

func NewWebsite(store repository.WebsiteRepository) Website {
	return Website{store: store}
}

func (service Website) List(ctx context.Context, userID string) ([]domain.Website, error) {
	return service.store.ListWebsites(ctx, userID)
}

func (service Website) Create(ctx context.Context, userID, name, domainName string) (domain.Website, error) {
	name, domainName = normalizeWebsiteInput(name, domainName)
	if name == "" {
		return domain.Website{}, ErrInvalidWebsiteName
	}
	return service.store.CreateWebsite(ctx, userID, name, domainName)
}

func (service Website) Get(ctx context.Context, userID, websiteID string) (domain.Website, error) {
	return service.store.GetWebsite(ctx, userID, websiteID)
}

func (service Website) Update(ctx context.Context, userID, websiteID, name, domainName string) (domain.Website, error) {
	name, domainName = normalizeWebsiteInput(name, domainName)
	if name == "" {
		return domain.Website{}, ErrInvalidWebsiteName
	}

	if err := service.store.UpdateWebsite(ctx, userID, websiteID, name, domainName); err != nil {
		return domain.Website{}, err
	}
	return service.store.GetWebsite(ctx, userID, websiteID)
}

func (service Website) Delete(ctx context.Context, userID, websiteID string) error {
	return service.store.DeleteWebsite(ctx, userID, websiteID)
}

func normalizeWebsiteInput(name, domain string) (string, string) {
	return strings.TrimSpace(name), strings.TrimSpace(domain)
}
