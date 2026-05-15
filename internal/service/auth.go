package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/model"
)

var ErrInvalidCredentials = errors.New("invalid username or password")

type UserStore interface {
	GetUserByUsername(ctx context.Context, username string) (model.User, error)
}

type Auth struct {
	users      UserStore
	isNotFound func(error) bool
}

func NewAuth(users UserStore, isNotFound func(error) bool) Auth {
	return Auth{users: users, isNotFound: isNotFound}
}

func (service Auth) Login(ctx context.Context, username, password string) (model.User, error) {
	user, err := service.users.GetUserByUsername(ctx, username)
	if err != nil {
		if service.isNotFound(err) {
			return model.User{}, ErrInvalidCredentials
		}
		return model.User{}, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return model.User{}, ErrInvalidCredentials
	}

	return user, nil
}
