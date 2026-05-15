package repository

import (
	"context"
	"fmt"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	storesqlc "github.com/kaixianzheng1216-creator/go-fetch/internal/repository/sqlc"

	"github.com/google/uuid"
)

func (store *Store) ListWebsites(ctx context.Context, userID string) ([]domain.Website, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("解析用户 ID 失败: %w", err)
	}

	rows, err := store.queries.ListWebsites(ctx, userUUID)
	if err != nil {
		return nil, fmt.Errorf("查询站点列表失败: %w", err)
	}

	websites := make([]domain.Website, 0, len(rows))
	for _, row := range rows {
		websites = append(websites, toWebsite(row.ID, row.Name, row.Domain, row.CreatedAt))
	}

	return websites, nil
}

func (store *Store) CreateWebsite(ctx context.Context, userID, name, domainName string) (domain.Website, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return domain.Website{}, fmt.Errorf("解析用户 ID 失败: %w", err)
	}

	row, err := store.queries.CreateWebsite(ctx, storesqlc.CreateWebsiteParams{
		ID:     uuid.New(),
		UserID: userUUID,
		Name:   name,
		Domain: domainName,
	})
	if err != nil {
		return domain.Website{}, fmt.Errorf("创建站点失败: %w", err)
	}

	return toWebsite(row.ID, row.Name, row.Domain, row.CreatedAt), nil
}

func (store *Store) GetWebsite(ctx context.Context, userID, websiteID string) (domain.Website, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return domain.Website{}, fmt.Errorf("解析用户 ID 失败: %w", err)
	}

	websiteUUID, err := uuid.Parse(websiteID)
	if err != nil {
		return domain.Website{}, fmt.Errorf("解析站点 ID 失败: %w", err)
	}

	row, err := store.queries.GetWebsite(ctx, storesqlc.GetWebsiteParams{ID: websiteUUID, UserID: userUUID})
	if err != nil {
		return domain.Website{}, fmt.Errorf("查询站点失败: %w", mapNotFound(err))
	}

	return toWebsite(row.ID, row.Name, row.Domain, row.CreatedAt), nil
}

func (store *Store) GetWebsiteForCollection(ctx context.Context, websiteID string) (domain.Website, error) {
	websiteUUID, err := uuid.Parse(websiteID)
	if err != nil {
		return domain.Website{}, fmt.Errorf("解析站点 ID 失败: %w", err)
	}

	row, err := store.queries.GetWebsiteForCollection(ctx, websiteUUID)
	if err != nil {
		return domain.Website{}, fmt.Errorf("查询采集站点失败: %w", mapNotFound(err))
	}

	return toWebsite(row.ID, row.Name, row.Domain, row.CreatedAt), nil
}

func (store *Store) UpdateWebsite(ctx context.Context, userID, websiteID, name, domainName string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("解析用户 ID 失败: %w", err)
	}

	websiteUUID, err := uuid.Parse(websiteID)
	if err != nil {
		return fmt.Errorf("解析站点 ID 失败: %w", err)
	}

	rows, err := store.queries.UpdateWebsite(ctx, storesqlc.UpdateWebsiteParams{
		ID:     websiteUUID,
		UserID: userUUID,
		Name:   name,
		Domain: domainName,
	})
	if err != nil {
		return fmt.Errorf("更新站点失败: %w", err)
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (store *Store) DeleteWebsite(ctx context.Context, userID, websiteID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return fmt.Errorf("解析用户 ID 失败: %w", err)
	}

	websiteUUID, err := uuid.Parse(websiteID)
	if err != nil {
		return fmt.Errorf("解析站点 ID 失败: %w", err)
	}

	rows, err := store.queries.DeleteWebsite(ctx, storesqlc.DeleteWebsiteParams{ID: websiteUUID, UserID: userUUID})
	if err != nil {
		return fmt.Errorf("删除站点失败: %w", err)
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}
