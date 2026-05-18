package mvp_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/domain"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/httpapi"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/repository"
	"github.com/kaixianzheng1216-creator/go-fetch/internal/service"
)

var (
	testUserID    = uuid.MustParse("2f931f32-9707-4ff1-8ca6-e851352d1451")
	testWebsiteID = uuid.MustParse("8a7e7a10-7b51-43ef-9e85-874df7dd5f8b")
)

func TestCollectPreflightUsesCORS(t *testing.T) {
	request := httptest.NewRequest(http.MethodOptions, "/api/collect", nil)
	request.Header.Set("Origin", "https://tracked.example")
	request.Header.Set("Access-Control-Request-Method", http.MethodPost)
	response := httptest.NewRecorder()

	httpapi.New(&repository.Store{}, nil, httpapi.Config{}).ServeHTTP(response, request)

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
			name:           "missing client info",
			collectionType: "event",
			expectedError:  service.ErrMissingClientInfo,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			store := &fakeTrackingStore{}
			collect := service.NewCollector(store)

			params := service.CollectionParams{
				Client: testClientInfo(testCase.request),
				Type:   domain.CollectionType(testCase.collectionType),
				Payload: domain.CollectPayload{
					WebsiteID: testWebsiteID,
					URL:       "https://example.com/",
				},
			}
			err := collect.Collect(context.Background(), params)

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
				_, err := stats.Summary(context.Background(), service.StatsParams{
					UserID:    testUserID,
					WebsiteID: testWebsiteID,
					Range:     testDateRange(startAt, endAt),
				})
				return err
			},
		},
		{
			name: "pageviews",
			call: func(stats service.Stats) error {
				_, err := stats.Pageviews(context.Background(), service.PageviewsParams{
					StatsParams: service.StatsParams{
						UserID:    testUserID,
						WebsiteID: testWebsiteID,
						Range:     testDateRange(startAt, endAt),
					},
					Unit: domain.DateUnitHour,
				})
				return err
			},
		},
		{
			name: "metrics",
			call: func(stats service.Stats) error {
				_, err := stats.Metrics(context.Background(), service.MetricsParams{
					StatsParams: service.StatsParams{
						UserID:    testUserID,
						WebsiteID: testWebsiteID,
						Range:     testDateRange(startAt, endAt),
					},
					Type:  domain.MetricTypePath,
					Limit: 10,
				})
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

			_, err := stats.Metrics(context.Background(), service.MetricsParams{
				StatsParams: service.StatsParams{
					UserID:    testUserID,
					WebsiteID: testWebsiteID,
				},
				Type:  domain.MetricTypePath,
				Limit: testCase.limit,
			})
			if err != nil {
				t.Fatalf("Metrics() error = %v", err)
			}
			if store.metricLimit != testCase.expected {
				t.Fatalf("metricLimit = %d, want %d", store.metricLimit, testCase.expected)
			}
		})
	}
}

func testDateRange(startAt, endAt int64) service.DateRange {
	start := time.UnixMilli(startAt).UTC()
	end := time.UnixMilli(endAt).UTC()
	return service.DateRange{StartAt: &start, EndAt: &end}
}

func testClientInfo(request *http.Request) service.ClientInfo {
	if request == nil {
		return service.ClientInfo{}
	}
	return service.ClientInfo{
		IP:        request.RemoteAddr,
		UserAgent: request.UserAgent(),
	}
}

type fakeTrackingStore struct {
	lookupCalls int
	saveCalls   int
}

func (store *fakeTrackingStore) GetWebsiteForCollection(_ context.Context, websiteID uuid.UUID) (domain.Website, error) {
	store.lookupCalls++
	return domain.Website{ID: websiteID, Name: "Example"}, nil
}

func (store *fakeTrackingStore) SaveEvent(_ context.Context, _ domain.EventRecord) error {
	store.saveCalls++
	return nil
}

type fakeStatsStore struct {
	accessCalls int
	metricLimit int
}

func (store *fakeStatsStore) GetWebsite(_ context.Context, _ uuid.UUID, websiteID uuid.UUID) (domain.Website, error) {
	store.accessCalls++
	return domain.Website{ID: websiteID, Name: "Example"}, nil
}

func (*fakeStatsStore) WebsiteStats(_ context.Context, _ uuid.UUID, _, _ time.Time) (domain.WebsiteStats, error) {
	return domain.WebsiteStats{}, nil
}

func (*fakeStatsStore) WebsitePageviews(_ context.Context, _ uuid.UUID, _, _ time.Time, _ domain.DateUnit) ([]domain.PageviewBucket, error) {
	return nil, nil
}

func (store *fakeStatsStore) WebsiteMetrics(_ context.Context, _ uuid.UUID, _, _ time.Time, _ domain.MetricType, limit int) ([]domain.Metric, error) {
	store.metricLimit = limit
	return nil, nil
}
