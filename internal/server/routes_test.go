package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alexedwards/scs/v2"
)

func testApp() *App {
	return &App{
		sessions: scs.New(),
	}
}

func performRequest(method, path, body string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	testApp().Routes().ServeHTTP(rec, req)
	return rec
}

func TestRoutesRequireAuthReturnsJSONUnauthorized(t *testing.T) {
	rec := performRequest(http.MethodGet, "/api/me", "")

	assertStatus(t, rec, http.StatusUnauthorized)
	assertContentType(t, rec, "application/problem+json")
	assertBodyContains(t, rec, `"detail":"未登录或登录已失效"`)
}

func TestRoutesInvalidJSONReturnsProblemDetails(t *testing.T) {
	rec := performRequest(http.MethodPost, "/api/login", "{")

	assertStatus(t, rec, http.StatusBadRequest)
	assertBodyContains(t, rec, `"status":400`)
}

func TestRoutesValidateRequestBodies(t *testing.T) {
	tests := []struct {
		name string
		path string
		body string
	}{
		{name: "login", path: "/api/login", body: `{}`},
		{name: "collect", path: "/api/collect", body: `{"payload":{}}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := performRequest(http.MethodPost, tt.path, tt.body)

			assertStatus(t, rec, http.StatusUnprocessableEntity)
			assertBodyContains(t, rec, `"detail":"validation failed"`)
		})
	}
}

func TestRoutesServeHumaOpenAPI(t *testing.T) {
	rec := performRequest(http.MethodGet, "/openapi.json", "")

	assertStatus(t, rec, http.StatusOK)
	assertContentType(t, rec, "application/openapi+json")

	body := rec.Body.String()
	if !strings.Contains(body, `"/api/websites/{websiteID}/metrics"`) {
		t.Fatalf("OpenAPI response missing metrics path")
	}
	if !strings.Contains(body, `"sessionCookie"`) {
		t.Fatalf("OpenAPI response missing session cookie security scheme")
	}
}

func TestOpenAPIJSONIsGeneratedFromServerRoutes(t *testing.T) {
	body, err := OpenAPIJSON()
	if err != nil {
		t.Fatalf("generate OpenAPI JSON: %v", err)
	}
	if !strings.Contains(string(body), `"operationId": "websiteMetrics"`) {
		t.Fatalf("OpenAPI output missing websiteMetrics operation")
	}
}

func assertStatus(t *testing.T, rec *httptest.ResponseRecorder, want int) {
	t.Helper()

	if rec.Code != want {
		t.Fatalf("status = %d, want %d", rec.Code, want)
	}
}

func assertContentType(t *testing.T, rec *httptest.ResponseRecorder, want string) {
	t.Helper()

	if got := rec.Header().Get("Content-Type"); got != want {
		t.Fatalf("Content-Type = %q, want %q", got, want)
	}
}

func assertBodyContains(t *testing.T, rec *httptest.ResponseRecorder, want string) {
	t.Helper()

	if body := rec.Body.String(); !strings.Contains(body, want) {
		t.Fatalf("body = %q, want substring %q", body, want)
	}
}
