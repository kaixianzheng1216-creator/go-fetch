package store

import (
	"context"
	"errors"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	storedb "github.com/kaixianzheng1216-creator/go-fetch/internal/store/db"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *Store) EnsureAdmin(ctx context.Context, username, password string) error {
	count, err := s.queries.CountUsers(ctx)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	if password == "" {
		return errors.New("ADMIN_PASSWORD is required when bootstrapping the first admin user")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.queries.CreateUser(ctx, storedb.CreateUserParams{
		ID:           uuid.New(),
		Username:     username,
		PasswordHash: string(hash),
	})
}

func (s *Store) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	row, err := s.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return domain.User{}, mapNotFound(err)
	}
	return toUser(row.ID, row.Username, row.PasswordHash, row.CreatedAt), nil
}

func (s *Store) GetUserByID(ctx context.Context, userID string) (domain.User, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return domain.User{}, err
	}
	row, err := s.queries.GetUserByID(ctx, userUUID)
	if err != nil {
		return domain.User{}, mapNotFound(err)
	}
	return toUser(row.ID, row.Username, row.PasswordHash, row.CreatedAt), nil
}
