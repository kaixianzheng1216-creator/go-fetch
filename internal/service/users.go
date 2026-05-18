package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
)

type UserRepository interface {
	CountUsers(ctx context.Context) (int64, error)
	CreateUser(ctx context.Context, user domain.User) error
	GetUserByID(ctx context.Context, userID uuid.UUID) (domain.User, error)
}

type UserService struct {
	repository UserRepository
}

func NewUserService(repository UserRepository) UserService {
	return UserService{repository: repository}
}

func (svc UserService) EnsureAdminUser(ctx context.Context, username, password string) error {
	count, err := svc.repository.CountUsers(ctx)
	if err != nil {
		return fmt.Errorf("count users: %w", err)
	}
	if count > 0 {
		return nil
	}

	username = strings.TrimSpace(username)
	if username == "" {
		return fmt.Errorf("admin username cannot be empty")
	}
	if password == "" {
		return fmt.Errorf("admin password cannot be empty")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash admin password: %w", err)
	}

	user := domain.User{
		ID:           uuid.New(),
		Username:     username,
		PasswordHash: string(hash),
	}
	if err := svc.repository.CreateUser(ctx, user); err != nil {
		return fmt.Errorf("create admin user: %w", err)
	}

	return nil
}

func (svc UserService) GetByID(ctx context.Context, userID uuid.UUID) (domain.User, error) {
	return svc.repository.GetUserByID(ctx, userID)
}
