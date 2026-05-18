package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	storesqlc "github.com/kaixianzheng1216-creator/go-fetch/internal/repository/sqlc"
)

func (store *Store) CountUsers(ctx context.Context) (int64, error) {
	count, err := store.queries.CountUsers(ctx)
	if err != nil {
		return 0, fmt.Errorf("count users: %w", err)
	}

	return count, nil
}

func (store *Store) CreateUser(ctx context.Context, user domain.User) error {
	if err := store.queries.CreateUser(ctx, storesqlc.CreateUserParams{
		ID:           user.ID,
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
	}); err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (store *Store) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	row, err := store.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user by username: %w", mapNotFound(err))
	}

	return toUser(row.ID, row.Username, row.PasswordHash, row.CreatedAt, row.UpdatedAt, row.DeletedAt), nil
}

func (store *Store) GetUserByID(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	row, err := store.queries.GetUserByID(ctx, userID)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user by ID: %w", mapNotFound(err))
	}

	return toUser(row.ID, row.Username, row.PasswordHash, row.CreatedAt, row.UpdatedAt, row.DeletedAt), nil
}
