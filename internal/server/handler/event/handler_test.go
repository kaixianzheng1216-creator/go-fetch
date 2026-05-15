package event

import (
	"context"
	"net/http"
	"testing"

	"github.com/danielgtaylor/huma/v2"

	eventdomain "github.com/kaixianzheng1216-creator/go-fetch/internal/domain/event"
	websitedomain "github.com/kaixianzheng1216-creator/go-fetch/internal/domain/website"
)

type fakeEventStore struct{}

func (fakeEventStore) GetWebsiteForCollection(context.Context, string) (websitedomain.Website, error) {
	return websitedomain.Website{}, nil
}

func (fakeEventStore) SaveEvent(context.Context, eventdomain.EventInput) error {
	return nil
}

func TestCollectRejectsUnsupportedType(t *testing.T) {
	handler := New(fakeEventStore{}, func(context.Context) *http.Request { return nil }, func(error) bool { return false })

	_, err := handler.Collect(context.Background(), &collectRequest{
		Body: CollectRequest{Type: CollectionType("unknown")},
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
