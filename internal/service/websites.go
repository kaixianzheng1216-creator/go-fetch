package service

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
)

var ErrInvalidWebsiteName = errors.New("website name cannot be empty")

// WebsiteRepository persists user-owned websites.
type WebsiteRepository interface {
	ListWebsites(ctx context.Context, userID uuid.UUID) ([]domain.Website, error)
	CreateWebsite(ctx context.Context, userID uuid.UUID, name, domain string) (domain.Website, error)
	GetWebsite(ctx context.Context, userID, websiteID uuid.UUID) (domain.Website, error)
	UpdateWebsite(ctx context.Context, userID, websiteID uuid.UUID, name, domain string) (domain.Website, error)
	DeleteWebsite(ctx context.Context, userID, websiteID uuid.UUID) error
}

// CreateWebsiteParams contains website fields used during creation.
type CreateWebsiteParams struct {
	Name   string
	Domain string
}

// UpdateWebsiteParams contains website fields used during updates.
type UpdateWebsiteParams struct {
	Name   string
	Domain string
}

// WebsiteService manages user-owned websites.
type WebsiteService struct {
	repository WebsiteRepository
}

// NewWebsiteService returns a website service.
func NewWebsiteService(repository WebsiteRepository) WebsiteService {
	return WebsiteService{repository: repository}
}

// List returns websites owned by a user.
func (svc WebsiteService) List(ctx context.Context, userID uuid.UUID) ([]domain.Website, error) {
	return svc.repository.ListWebsites(ctx, userID)
}

// Create creates a user-owned website.
func (svc WebsiteService) Create(ctx context.Context, userID uuid.UUID, params CreateWebsiteParams) (domain.Website, error) {
	params.Name, params.Domain = normalizeWebsiteFields(params.Name, params.Domain)
	if params.Name == "" {
		return domain.Website{}, ErrInvalidWebsiteName
	}
	return svc.repository.CreateWebsite(ctx, userID, params.Name, params.Domain)
}

// Get returns a user-owned website.
func (svc WebsiteService) Get(ctx context.Context, userID, websiteID uuid.UUID) (domain.Website, error) {
	return svc.repository.GetWebsite(ctx, userID, websiteID)
}

// Update updates a user-owned website.
func (svc WebsiteService) Update(ctx context.Context, userID, websiteID uuid.UUID, params UpdateWebsiteParams) (domain.Website, error) {
	params.Name, params.Domain = normalizeWebsiteFields(params.Name, params.Domain)
	if params.Name == "" {
		return domain.Website{}, ErrInvalidWebsiteName
	}

	return svc.repository.UpdateWebsite(ctx, userID, websiteID, params.Name, params.Domain)
}

// Delete deletes a user-owned website.
func (svc WebsiteService) Delete(ctx context.Context, userID, websiteID uuid.UUID) error {
	return svc.repository.DeleteWebsite(ctx, userID, websiteID)
}

func normalizeWebsiteFields(name, domain string) (string, string) {
	return strings.TrimSpace(name), strings.TrimSpace(domain)
}
