package store

import (
	"context"
	"fmt"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	storedb "github.com/kaixianzheng1216-creator/go-fetch/internal/store/db"

	"github.com/google/uuid"
)

func (s *Store) ListWebsites(ctx context.Context, userID string) ([]domain.Website, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("parse user id: %w", err)
	}

	rows, err := s.queries.ListWebsites(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("list websites: %w", err)
	}

	websites := make([]domain.Website, 0, len(rows))
	for _, row := range rows {
		websites = append(websites, toWebsite(row.ID, row.Name, row.Domain, row.CreatedAt))
	}

	return websites, nil
}

func (s *Store) CreateWebsite(ctx context.Context, userID, name, websiteDomain string) (domain.Website, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return domain.Website{}, fmt.Errorf("parse user id: %w", err)
	}

	row, err := s.queries.CreateWebsite(ctx, storedb.CreateWebsiteParams{
		ID:     uuid.New(),
		UserID: userUUID,
		Name:   name,
		Domain: websiteDomain,
	})
	if err != nil {
		return domain.Website{}, fmt.Errorf("create website: %w", err)
	}

	return toWebsite(row.ID, row.Name, row.Domain, row.CreatedAt), nil
}

func (s *Store) GetWebsite(ctx context.Context, userID, websiteID string) (domain.Website, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return domain.Website{}, fmt.Errorf("parse user id: %w", err)
	}

	websiteUUID, err := uuid.Parse(websiteID)
	if err != nil {
		return domain.Website{}, fmt.Errorf("parse website id: %w", err)
	}

	row, err := s.queries.GetWebsite(ctx, storedb.GetWebsiteParams{ID: websiteUUID, UserID: userUUID})
	if err != nil {
		return domain.Website{}, fmt.Errorf("get website: %w", mapNotFound(err))
	}

	return toWebsite(row.ID, row.Name, row.Domain, row.CreatedAt), nil
}

func (s *Store) GetWebsiteForCollection(ctx context.Context, websiteID string) (domain.Website, error) {
	websiteUUID, err := uuid.Parse(websiteID)
	if err != nil {
		return domain.Website{}, fmt.Errorf("parse website id: %w", err)
	}

	row, err := s.queries.GetWebsiteForCollection(ctx, websiteUUID)
	if err != nil {
		return domain.Website{}, fmt.Errorf("get website for collection: %w", mapNotFound(err))
	}

	return toWebsite(row.ID, row.Name, row.Domain, row.CreatedAt), nil
}

func (s *Store) UpdateWebsite(ctx context.Context, userID, websiteID, name, websiteDomain string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("parse user id: %w", err)
	}

	websiteUUID, err := uuid.Parse(websiteID)
	if err != nil {
		return fmt.Errorf("parse website id: %w", err)
	}

	rows, err := s.queries.UpdateWebsite(ctx, storedb.UpdateWebsiteParams{
		ID:     websiteUUID,
		UserID: userUUID,
		Name:   name,
		Domain: websiteDomain,
	})
	if err != nil {
		return fmt.Errorf("update website: %w", err)
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *Store) DeleteWebsite(ctx context.Context, userID, websiteID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("parse user id: %w", err)
	}

	websiteUUID, err := uuid.Parse(websiteID)
	if err != nil {
		return fmt.Errorf("parse website id: %w", err)
	}

	rows, err := s.queries.DeleteWebsite(ctx, storedb.DeleteWebsiteParams{ID: websiteUUID, UserID: userUUID})
	if err != nil {
		return fmt.Errorf("delete website: %w", err)
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}
