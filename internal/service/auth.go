package service

import (
	"context"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
)

var ErrInvalidCredentials = errors.New("invalid username or password")

// UserRepository provides user lookups for authentication.
type UserRepository interface {
	GetUserByUsername(ctx context.Context, username string) (domain.User, error)
}

type Auth struct {
	users UserRepository
}

func NewAuth(users UserRepository) Auth {
	return Auth{users: users}
}

func (service Auth) Login(ctx context.Context, username, password string) (domain.User, error) {
	username = strings.TrimSpace(username)
	if username == "" || password == "" {
		return domain.User{}, ErrInvalidCredentials
	}

	user, err := service.users.GetUserByUsername(ctx, username)
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
