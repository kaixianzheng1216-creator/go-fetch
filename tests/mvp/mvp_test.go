package mvp_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/repository"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/router"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/service"
)

func TestCollectPreflightUsesCORS(t *testing.T) {
	request := httptest.NewRequest(http.MethodOptions, "/api/collect", nil)
	request.Header.Set("Origin", "https://tracked.example")
	request.Header.Set("Access-Control-Request-Method", http.MethodPost)
	response := httptest.NewRecorder()

	router.New(&repository.Store{}).ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if actual := response.Header().Get("Access-Control-Allow-Origin"); actual != "*" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want *", actual)
	}
}

func TestCollectErrors(t *testing.T) {
	tests := []struct {
		name             string
		collectionType   string
		request          *http.Request
		expectedError    error
		expectedLookup   int
		expectedSaveCall int
	}{
		{
			name:           "unsupported collection type",
			collectionType: "pageview",
			expectedError:  service.ErrUnsupportedCollectionType,
		},
		{
			name:           "missing http request after website lookup",
			collectionType: "event",
			expectedError:  service.ErrMissingHTTPRequest,
			expectedLookup: 1,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			store := &fakeTrackingStore{}
			collect := service.NewCollect(store)

			err := collect.Collect(context.Background(), testCase.request, testCase.collectionType, domain.CollectPayload{
				WebsiteID: "8a7e7a10-7b51-43ef-9e85-874df7dd5f8b",
				URL:       "https://example.com/",
			})

			if !errors.Is(err, testCase.expectedError) {
				t.Fatalf("error = %v, want %v", err, testCase.expectedError)
			}
			if store.lookupCalls != testCase.expectedLookup {
				t.Fatalf("lookupCalls = %d, want %d", store.lookupCalls, testCase.expectedLookup)
			}
			if store.saveCalls != testCase.expectedSaveCall {
				t.Fatalf("saveCalls = %d, want %d", store.saveCalls, testCase.expectedSaveCall)
			}
		})
	}
}

func TestStatsRejectsInvalidDateRange(t *testing.T) {
	startAt := time.Date(2026, 5, 15, 10, 0, 0, 0, time.UTC).UnixMilli()
	endAt := time.Date(2026, 5, 15, 9, 0, 0, 0, time.UTC).UnixMilli()

	tests := []struct {
		name string
		call func(service.Stats) error
	}{
		{
			name: "stats",
			call: func(stats service.Stats) error {
				_, err := stats.WebsiteStats(context.Background(), "user-1", "website-1", &startAt, &endAt)
				return err
			},
		},
		{
			name: "pageviews",
			call: func(stats service.Stats) error {
				_, err := stats.WebsitePageviews(context.Background(), "user-1", "website-1", &startAt, &endAt, "hour")
				return err
			},
		},
		{
			name: "metrics",
			call: func(stats service.Stats) error {
				_, err := stats.WebsiteMetrics(context.Background(), "user-1", "website-1", &startAt, &endAt, "path", 10)
				return err
			},
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			store := &fakeStatsStore{}
			err := testCase.call(service.NewStats(store))

			if !errors.Is(err, service.ErrInvalidDateRange) {
				t.Fatalf("error = %v, want %v", err, service.ErrInvalidDateRange)
			}
			if store.accessCalls != 0 {
				t.Fatalf("accessCalls = %d, want 0", store.accessCalls)
			}
		})
	}
}

func TestStatsWebsiteMetricsNormalizesLimit(t *testing.T) {
	tests := []struct {
		name     string
		limit    int
		expected int
	}{
		{name: "default", limit: 0, expected: domain.DefaultMetricLimit},
		{name: "too high", limit: domain.MaxMetricLimit + 1, expected: domain.DefaultMetricLimit},
		{name: "valid", limit: 25, expected: 25},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			store := &fakeStatsStore{}
			stats := service.NewStats(store)

			_, err := stats.WebsiteMetrics(context.Background(), "user-1", "website-1", nil, nil, "path", testCase.limit)
			if err != nil {
				t.Fatalf("WebsiteMetrics() error = %v", err)
			}
			if store.metricLimit != testCase.expected {
				t.Fatalf("metricLimit = %d, want %d", store.metricLimit, testCase.expected)
			}
		})
	}
}

type fakeTrackingStore struct {
	lookupCalls int
	saveCalls   int
}

func (store *fakeTrackingStore) GetWebsiteForCollection(_ context.Context, websiteID string) (domain.Website, error) {
	store.lookupCalls++
	return domain.Website{ID: websiteID, Name: "Example"}, nil
}

func (store *fakeTrackingStore) SaveEvent(_ context.Context, _ domain.EventInput) error {
	store.saveCalls++
	return nil
}

type fakeStatsStore struct {
	accessCalls int
	metricLimit int
}

func (store *fakeStatsStore) GetWebsite(_ context.Context, _ string, websiteID string) (domain.Website, error) {
	store.accessCalls++
	return domain.Website{ID: websiteID, Name: "Example"}, nil
}

func (*fakeStatsStore) WebsiteStats(_ context.Context, _ string, _, _ time.Time) (domain.WebsiteStats, error) {
	return domain.WebsiteStats{}, nil
}

func (*fakeStatsStore) WebsitePageviews(_ context.Context, _ string, _, _ time.Time, _ domain.DateUnit) ([]domain.PageviewPoint, error) {
	return nil, nil
}

func (store *fakeStatsStore) WebsiteMetrics(_ context.Context, _ string, _, _ time.Time, _ domain.MetricType, limit int) ([]domain.MetricRow, error) {
	store.metricLimit = limit
	return nil, nil
}
