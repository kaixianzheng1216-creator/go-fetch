package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	storesqlc "github.com/kaixianzheng1216-creator/go-fetch/internal/repository/sqlc"
)

func (store *Store) ListWebsites(ctx context.Context, userID uuid.UUID) ([]domain.Website, error) {
	rows, err := store.queries.ListWebsites(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list websites: %w", err)
	}

	websites := make([]domain.Website, len(rows))
	for i, row := range rows {
		websites[i] = toWebsite(row.ID, row.Name, row.Domain, row.CreatedAt)
	}

	return websites, nil
}

func (store *Store) CreateWebsite(ctx context.Context, userID uuid.UUID, name, siteDomain string) (domain.Website, error) {
	row, err := store.queries.CreateWebsite(ctx, storesqlc.CreateWebsiteParams{
		ID:     uuid.New(),
		UserID: userID,
		Name:   name,
		Domain: siteDomain,
	})
	if err != nil {
		return domain.Website{}, fmt.Errorf("create website: %w", err)
	}

	return toWebsite(row.ID, row.Name, row.Domain, row.CreatedAt), nil
}

func (store *Store) GetWebsite(ctx context.Context, userID, websiteID uuid.UUID) (domain.Website, error) {
	row, err := store.queries.GetWebsite(ctx, storesqlc.GetWebsiteParams{ID: websiteID, UserID: userID})
	if err != nil {
		return domain.Website{}, fmt.Errorf("get website: %w", mapNotFound(err))
	}

	return toWebsite(row.ID, row.Name, row.Domain, row.CreatedAt), nil
}

func (store *Store) GetWebsiteForCollection(ctx context.Context, websiteID uuid.UUID) (domain.Website, error) {
	row, err := store.queries.GetWebsiteForCollection(ctx, websiteID)
	if err != nil {
		return domain.Website{}, fmt.Errorf("get collection website: %w", mapNotFound(err))
	}

	return toWebsite(row.ID, row.Name, row.Domain, row.CreatedAt), nil
}

func (store *Store) UpdateWebsite(ctx context.Context, userID, websiteID uuid.UUID, name, siteDomain string) (domain.Website, error) {
	row, err := store.queries.UpdateWebsite(ctx, storesqlc.UpdateWebsiteParams{
		ID:     websiteID,
		UserID: userID,
		Name:   name,
		Domain: siteDomain,
	})
	if err != nil {
		return domain.Website{}, fmt.Errorf("update website: %w", mapNotFound(err))
	}

	return toWebsite(row.ID, row.Name, row.Domain, row.CreatedAt), nil
}

func (store *Store) DeleteWebsite(ctx context.Context, userID, websiteID uuid.UUID) error {
	rows, err := store.queries.DeleteWebsite(ctx, storesqlc.DeleteWebsiteParams{ID: websiteID, UserID: userID})
	if err != nil {
		return fmt.Errorf("delete website: %w", err)
	}

	if rows == 0 {
		return domain.ErrNotFound
	}

	return nil
}
