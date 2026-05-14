package store

import (
	"context"
	"fmt"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	storedb "github.com/kaixianzheng1216-creator/go-fetch/internal/store/db"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *Store) EnsureAdmin(ctx context.Context, username, password string) error {
	count, err := s.queries.CountUsers(ctx)

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

	if err := s.queries.CreateUser(ctx, storedb.CreateUserParams{
		ID:           uuid.New(),
		Username:     username,
		PasswordHash: string(hash),
		DisplayName:  username,
	}); err != nil {
		return fmt.Errorf("create admin user: %w", err)
	}

	return nil
}

func (s *Store) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	row, err := s.queries.GetUserByUsername(ctx, username)

	if err != nil {
		return domain.User{}, fmt.Errorf("get user by username: %w", mapNotFound(err))
	}

	return toUser(row.ID, row.Username, row.PasswordHash, row.LogoUrl, row.DisplayName, row.CreatedAt, row.UpdatedAt, row.DeletedAt), nil
}

func (s *Store) GetUserByID(ctx context.Context, userID string) (domain.User, error) {
	userUUID, err := uuid.Parse(userID)

	if err != nil {
		return domain.User{}, fmt.Errorf("parse user id: %w", err)
	}

	row, err := s.queries.GetUserByID(ctx, userUUID)

	if err != nil {
		return domain.User{}, fmt.Errorf("get user by id: %w", mapNotFound(err))
	}

	return toUser(row.ID, row.Username, row.PasswordHash, row.LogoUrl, row.DisplayName, row.CreatedAt, row.UpdatedAt, row.DeletedAt), nil
}
