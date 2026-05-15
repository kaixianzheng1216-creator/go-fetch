package user

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2"

	userdomain "github.com/kaixianzheng1216-creator/go-fetch/internal/domain/user"
)

var errUserNotFound = errors.New("user not found")

type fakeUserStore struct{}

func (fakeUserStore) GetUserByUsername(context.Context, string) (userdomain.User, error) {
	return userdomain.User{}, errUserNotFound
}

type fakeSessions struct{}

func (fakeSessions) RenewToken(context.Context) error {
	return nil
}

func (fakeSessions) Put(context.Context, string, any) {}

func (fakeSessions) Destroy(context.Context) error {
	return nil
}

func TestLoginRejectsUnknownUser(t *testing.T) {
	handler := New(
		fakeUserStore{},
		fakeSessions{},
		"user_id",
		func(context.Context) userdomain.User { return userdomain.User{} },
		func(err error) bool { return errors.Is(err, errUserNotFound) },
	)

	_, err := handler.Login(context.Background(), &loginRequest{
		Body: LoginRequest{Username: "missing", Password: "secret"},
	})

	assertStatusError(t, err, http.StatusUnauthorized)
}

func assertStatusError(t *testing.T, err error, want int) {
	t.Helper()

	statusErr, ok := err.(huma.StatusError)
	if !ok {
		t.Fatalf("error = %T, want huma.StatusError", err)
	}
	if statusErr.GetStatus() != want {
		t.Fatalf("status = %d, want %d", statusErr.GetStatus(), want)
	}
}
