package summary

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/danielgtaylor/huma/v2"

	eventdomain "github.com/kaixianzheng1216-creator/go-fetch/internal/domain/event"
	userdomain "github.com/kaixianzheng1216-creator/go-fetch/internal/domain/user"
	websitedomain "github.com/kaixianzheng1216-creator/go-fetch/internal/domain/website"
)

type fakeSummaryStore struct{}

func (fakeSummaryStore) GetWebsite(context.Context, string, string) (websitedomain.Website, error) {
	return websitedomain.Website{}, nil
}

func (fakeSummaryStore) WebsiteStats(context.Context, string, time.Time, time.Time) (eventdomain.WebsiteStats, error) {
	return eventdomain.WebsiteStats{}, nil
}

func (fakeSummaryStore) Pageviews(context.Context, string, time.Time, time.Time, eventdomain.DateUnit) ([]eventdomain.PageviewPoint, error) {
	return nil, nil
}

func (fakeSummaryStore) Metrics(context.Context, string, time.Time, time.Time, eventdomain.MetricType, int) ([]eventdomain.MetricRow, error) {
	return nil, nil
}

func TestMetricsRejectsUnsupportedType(t *testing.T) {
	handler := New(
		fakeSummaryStore{},
		func(context.Context) userdomain.User { return userdomain.User{ID: "user-id"} },
		func(err error) error { return err },
	)

	_, err := handler.Metrics(context.Background(), &metricsRequest{
		WebsiteID: "website-id",
		Type:      MetricTypeParam("unknown"),
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
