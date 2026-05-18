// Package service contains application business operations.
package service

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
)

// ErrInvalidCredentials indicates an invalid login attempt.
var ErrInvalidCredentials = errors.New("invalid username or password")

// AuthUserRepository provides user lookups for authentication.
type AuthUserRepository interface {
	GetUserByUsername(ctx context.Context, username string) (domain.User, error)
}

type AuthService struct {
	users AuthUserRepository
}

func NewAuthService(users AuthUserRepository) AuthService {
	return AuthService{users: users}
}

func (svc AuthService) Login(ctx context.Context, username, password string) (domain.User, error) {
	username = strings.TrimSpace(username)
	if username == "" || password == "" {
		return domain.User{}, ErrInvalidCredentials
	}

	user, err := svc.users.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return domain.User{}, ErrInvalidCredentials
		}
		return domain.User{}, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return domain.User{}, ErrInvalidCredentials
	}

	return user, nil
}
