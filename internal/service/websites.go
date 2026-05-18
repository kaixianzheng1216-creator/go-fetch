package service

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
)

var ErrInvalidWebsiteName = errors.New("website name cannot be empty")

type WebsiteRepository interface {
	ListWebsites(ctx context.Context, userID uuid.UUID) ([]domain.Website, error)
	CreateWebsite(ctx context.Context, userID uuid.UUID, name, domainName string) (domain.Website, error)
	GetWebsite(ctx context.Context, userID, websiteID uuid.UUID) (domain.Website, error)
	UpdateWebsite(ctx context.Context, userID, websiteID uuid.UUID, name, domainName string) (domain.Website, error)
	DeleteWebsite(ctx context.Context, userID, websiteID uuid.UUID) error
}

type WebsiteInput struct {
	Name   string
	Domain string
}

type WebsiteService struct {
	repository WebsiteRepository
}

func NewWebsiteService(repository WebsiteRepository) WebsiteService {
	return WebsiteService{repository: repository}
}

func (svc WebsiteService) List(ctx context.Context, userID uuid.UUID) ([]domain.Website, error) {
	return svc.repository.ListWebsites(ctx, userID)
}

func (svc WebsiteService) Create(ctx context.Context, userID uuid.UUID, input WebsiteInput) (domain.Website, error) {
	input, err := normalizeWebsiteInput(input)
	if err != nil {
		return domain.Website{}, err
	}

	return svc.repository.CreateWebsite(ctx, userID, input.Name, input.Domain)
}

func (svc WebsiteService) Get(ctx context.Context, userID, websiteID uuid.UUID) (domain.Website, error) {
	return svc.repository.GetWebsite(ctx, userID, websiteID)
}

func (svc WebsiteService) Update(ctx context.Context, userID, websiteID uuid.UUID, input WebsiteInput) (domain.Website, error) {
	input, err := normalizeWebsiteInput(input)
	if err != nil {
		return domain.Website{}, err
	}

	return svc.repository.UpdateWebsite(ctx, userID, websiteID, input.Name, input.Domain)
}

func (svc WebsiteService) Delete(ctx context.Context, userID, websiteID uuid.UUID) error {
	return svc.repository.DeleteWebsite(ctx, userID, websiteID)
}

func normalizeWebsiteInput(input WebsiteInput) (WebsiteInput, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Domain = strings.TrimSpace(input.Domain)
	if input.Name == "" {
		return WebsiteInput{}, ErrInvalidWebsiteName
	}

	return input, nil
}
