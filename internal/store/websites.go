package store

import (
	"context"

	"go-fetch/internal/domain"
	storedb "go-fetch/internal/store/db"

	"github.com/google/uuid"
)

func (s *Store) ListWebsites(ctx context.Context, userID string) ([]domain.Website, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	rows, err := s.queries.ListWebsites(ctx, userUUID)
	if err != nil {
		return nil, err
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
		return domain.Website{}, err
	}
	row, err := s.queries.CreateWebsite(ctx, storedb.CreateWebsiteParams{
		ID:     uuid.New(),
		UserID: userUUID,
		Name:   name,
		Domain: websiteDomain,
	})
	if err != nil {
		return domain.Website{}, err
	}
	return toWebsite(row.ID, row.Name, row.Domain, row.CreatedAt), nil
}

func (s *Store) GetWebsite(ctx context.Context, userID, websiteID string) (domain.Website, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return domain.Website{}, err
	}
	websiteUUID, err := uuid.Parse(websiteID)
	if err != nil {
		return domain.Website{}, err
	}
	row, err := s.queries.GetWebsite(ctx, storedb.GetWebsiteParams{ID: websiteUUID, UserID: userUUID})
	if err != nil {
		return domain.Website{}, mapNotFound(err)
	}
	return toWebsite(row.ID, row.Name, row.Domain, row.CreatedAt), nil
}

func (s *Store) GetWebsiteForCollection(ctx context.Context, websiteID string) (domain.Website, error) {
	websiteUUID, err := uuid.Parse(websiteID)
	if err != nil {
		return domain.Website{}, err
	}
	row, err := s.queries.GetWebsiteForCollection(ctx, websiteUUID)
	if err != nil {
		return domain.Website{}, mapNotFound(err)
	}
	return toWebsite(row.ID, row.Name, row.Domain, row.CreatedAt), nil
}

func (s *Store) UpdateWebsite(ctx context.Context, userID, websiteID, name, websiteDomain string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	websiteUUID, err := uuid.Parse(websiteID)
	if err != nil {
		return err
	}
	rows, err := s.queries.UpdateWebsite(ctx, storedb.UpdateWebsiteParams{
		ID:     websiteUUID,
		UserID: userUUID,
		Name:   name,
		Domain: websiteDomain,
	})
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) DeleteWebsite(ctx context.Context, userID, websiteID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	websiteUUID, err := uuid.Parse(websiteID)
	if err != nil {
		return err
	}
	rows, err := s.queries.DeleteWebsite(ctx, storedb.DeleteWebsiteParams{ID: websiteUUID, UserID: userUUID})
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
