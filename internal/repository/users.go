package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	storesqlc "github.com/kaixianzheng1216-creator/go-fetch/internal/repository/sqlc"
)

func (store *Store) EnsureAdminUser(ctx context.Context, username, password string) error {
	count, err := store.queries.CountUsers(ctx)
	if err != nil {
		return fmt.Errorf("count users: %w", err)
	}
	if count > 0 {
		return nil
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash admin password: %w", err)
	}

	if err := store.queries.CreateUser(ctx, storesqlc.CreateUserParams{
		ID:           uuid.New(),
		Username:     username,
		PasswordHash: string(hash),
	}); err != nil {
		return fmt.Errorf("create admin user: %w", err)
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

func (store *Store) GetUserByID(ctx context.Context, userID string) (domain.User, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return domain.User{}, fmt.Errorf("parse user ID: %w", err)
	}

	row, err := store.queries.GetUserByID(ctx, userUUID)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user by ID: %w", mapNotFound(err))
	}

	return toUser(row.ID, row.Username, row.PasswordHash, row.CreatedAt, row.UpdatedAt, row.DeletedAt), nil
}
