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
	UpdateWebsite(ctx context.Context, userID, websiteID uuid.UUID, name, domainName string) error
	DeleteWebsite(ctx context.Context, userID, websiteID uuid.UUID) error
}

// WebsiteInput contains editable website fields.
type WebsiteInput struct {
	Name   string
	Domain string
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

func (service Websites) Create(ctx context.Context, userID uuid.UUID, input WebsiteInput) (domain.Website, error) {
	input = normalizeWebsiteInput(input)
	if input.Name == "" {
		return domain.Website{}, ErrInvalidWebsiteName
	}
	return service.repository.CreateWebsite(ctx, userID, input.Name, input.Domain)
}

func (service Websites) Get(ctx context.Context, userID, websiteID uuid.UUID) (domain.Website, error) {
	return service.repository.GetWebsite(ctx, userID, websiteID)
}

func (service Websites) Update(ctx context.Context, userID, websiteID uuid.UUID, input WebsiteInput) (domain.Website, error) {
	input = normalizeWebsiteInput(input)
	if input.Name == "" {
		return domain.Website{}, ErrInvalidWebsiteName
	}

	if err := service.repository.UpdateWebsite(ctx, userID, websiteID, input.Name, input.Domain); err != nil {
		return domain.Website{}, err
	}
	return service.repository.GetWebsite(ctx, userID, websiteID)
}

func (service Websites) Delete(ctx context.Context, userID, websiteID uuid.UUID) error {
	return service.repository.DeleteWebsite(ctx, userID, websiteID)
}

func normalizeWebsiteInput(input WebsiteInput) WebsiteInput {
	input.Name = strings.TrimSpace(input.Name)
	input.Domain = strings.TrimSpace(input.Domain)
	return input
}
