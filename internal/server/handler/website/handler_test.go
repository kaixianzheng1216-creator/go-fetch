package website

import (
	"context"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2"

	userdomain "github.com/kaixianzheng1216-creator/go-fetch/internal/domain/user"
	websitedomain "github.com/kaixianzheng1216-creator/go-fetch/internal/domain/website"
)

type fakeWebsiteStore struct{}

func (fakeWebsiteStore) ListWebsites(context.Context, string) ([]websitedomain.Website, error) {
	return nil, nil
}

func (fakeWebsiteStore) CreateWebsite(context.Context, string, string, string) (websitedomain.Website, error) {
	return websitedomain.Website{}, nil
}

func (fakeWebsiteStore) GetWebsite(context.Context, string, string) (websitedomain.Website, error) {
	return websitedomain.Website{}, nil
}

func (fakeWebsiteStore) UpdateWebsite(context.Context, string, string, string, string) error {
	return nil
}

func (fakeWebsiteStore) DeleteWebsite(context.Context, string, string) error {
	return nil
}

func TestCreateRejectsBlankName(t *testing.T) {
	handler := New(
		fakeWebsiteStore{},
		func(context.Context) userdomain.User { return userdomain.User{ID: "user-id"} },
		func(err error) error { return err },
	)

	_, err := handler.Create(context.Background(), &websiteRequest{
		Body: WebsiteRequest{Name: "   "},
	})

	assertStatusError(t, err, http.StatusBadRequest)
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
