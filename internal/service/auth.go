package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
)

var ErrInvalidCredentials = errors.New("invalid username or password")

type AuthStore interface {
	GetUserByUsername(ctx context.Context, username string) (domain.User, error)
}

type Auth struct {
	users      AuthStore
	isNotFound func(error) bool
}

func NewAuth(users AuthStore, isNotFound func(error) bool) Auth {
	return Auth{users: users, isNotFound: isNotFound}
}

func (service Auth) Login(ctx context.Context, username, password string) (domain.User, error) {
	user, err := service.users.GetUserByUsername(ctx, username)
	if err != nil {
		if service.isNotFound(err) {
			return domain.User{}, ErrInvalidCredentials
		}
		return domain.User{}, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return domain.User{}, ErrInvalidCredentials
	}

	return user, nil
}
