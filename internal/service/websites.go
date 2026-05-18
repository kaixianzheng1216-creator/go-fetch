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
	CreateWebsite(ctx context.Context, userID uuid.UUID, name, domainName string) (domain.Website, error)
	GetWebsite(ctx context.Context, userID, websiteID uuid.UUID) (domain.Website, error)
	UpdateWebsite(ctx context.Context, userID, websiteID uuid.UUID, name, domainName string) (domain.Website, error)
	DeleteWebsite(ctx context.Context, userID, websiteID uuid.UUID) error
}

// WebsiteParams contains editable website fields.
type WebsiteParams struct {
	Name       string
	DomainName string
}

type Websites struct {
	repository WebsiteRepository
}

func NewWebsites(repository WebsiteRepository) Websites {
	return Websites{repository: repository}
}

func (service Websites) List(ctx context.Context, userID uuid.UUID) ([]domain.Website, error) {
	return service.repository.ListWebsites(ctx, userID)
}

func (service Websites) Create(ctx context.Context, userID uuid.UUID, params WebsiteParams) (domain.Website, error) {
	params = normalizeWebsiteParams(params)
	if params.Name == "" {
		return domain.Website{}, ErrInvalidWebsiteName
	}
	return service.repository.CreateWebsite(ctx, userID, params.Name, params.DomainName)
}

func (service Websites) Get(ctx context.Context, userID, websiteID uuid.UUID) (domain.Website, error) {
	return service.repository.GetWebsite(ctx, userID, websiteID)
}

func (service Websites) Update(ctx context.Context, userID, websiteID uuid.UUID, params WebsiteParams) (domain.Website, error) {
	params = normalizeWebsiteParams(params)
	if params.Name == "" {
		return domain.Website{}, ErrInvalidWebsiteName
	}

	return service.repository.UpdateWebsite(ctx, userID, websiteID, params.Name, params.DomainName)
}

func (service Websites) Delete(ctx context.Context, userID, websiteID uuid.UUID) error {
	return service.repository.DeleteWebsite(ctx, userID, websiteID)
}

func normalizeWebsiteParams(params WebsiteParams) WebsiteParams {
	params.Name = strings.TrimSpace(params.Name)
	params.DomainName = strings.TrimSpace(params.DomainName)
	return params
}
