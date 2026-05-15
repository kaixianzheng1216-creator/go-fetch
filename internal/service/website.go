package service

import (
	"context"
	"errors"
	"strings"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/model"
)

var ErrInvalidWebsiteName = errors.New("website name cannot be empty")

type WebsiteStore interface {
	ListWebsites(ctx context.Context, userID string) ([]model.Website, error)
	CreateWebsite(ctx context.Context, userID, name, domainName string) (model.Website, error)
	GetWebsite(ctx context.Context, userID, websiteID string) (model.Website, error)
	UpdateWebsite(ctx context.Context, userID, websiteID, name, domainName string) error
	DeleteWebsite(ctx context.Context, userID, websiteID string) error
}

type Website struct {
	store WebsiteStore
}

func NewWebsite(store WebsiteStore) Website {
	return Website{store: store}
}

func (service Website) List(ctx context.Context, userID string) ([]model.Website, error) {
	return service.store.ListWebsites(ctx, userID)
}

func (service Website) Create(ctx context.Context, userID, name, domain string) (model.Website, error) {
	name, domain = normalizeWebsiteInput(name, domain)
	if name == "" {
		return model.Website{}, ErrInvalidWebsiteName
	}
	return service.store.CreateWebsite(ctx, userID, name, domain)
}

func (service Website) Get(ctx context.Context, userID, websiteID string) (model.Website, error) {
	return service.store.GetWebsite(ctx, userID, websiteID)
}

func (service Website) Update(ctx context.Context, userID, websiteID, name, domain string) (model.Website, error) {
	name, domain = normalizeWebsiteInput(name, domain)
	if name == "" {
		return model.Website{}, ErrInvalidWebsiteName
	}

	if err := service.store.UpdateWebsite(ctx, userID, websiteID, name, domain); err != nil {
		return model.Website{}, err
	}
	return service.store.GetWebsite(ctx, userID, websiteID)
}

func (service Website) Delete(ctx context.Context, userID, websiteID string) error {
	return service.store.DeleteWebsite(ctx, userID, websiteID)
}

func normalizeWebsiteInput(name, domain string) (string, string) {
	return strings.TrimSpace(name), strings.TrimSpace(domain)
}
