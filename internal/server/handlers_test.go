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

func TestRoutesRequireAuthReturnsJSONUnauthorized(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
	testApp().Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d", rec.Code)
	}
	if contentType := rec.Header().Get("Content-Type"); contentType != "application/problem+json" {
		t.Fatalf("Content-Type = %q", contentType)
	}
	if body := rec.Body.String(); !strings.Contains(body, `"detail":"unauthorized"`) {
		t.Fatalf("body = %q", body)
	}
}

func TestRoutesInvalidJSONReturnsProblemDetails(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader("{"))
	req.Header.Set("Content-Type", "application/json")
	testApp().Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d", rec.Code)
	}
	if body := rec.Body.String(); !strings.Contains(body, `"detail":"validation failed"`) {
		t.Fatalf("body = %q", body)
	}
}

func TestRoutesServeHumaOpenAPI(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/openapi.json", nil)
	testApp().Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if contentType := rec.Header().Get("Content-Type"); contentType != "application/openapi+json" {
		t.Fatalf("Content-Type = %q", contentType)
	}

	body := rec.Body.String()
	if !strings.Contains(body, `"/api/websites/{websiteID}/metrics"`) {
		t.Fatalf("OpenAPI body does not include metrics path")
	}
	if !strings.Contains(body, `"sessionCookie"`) {
		t.Fatalf("OpenAPI body does not include session cookie security scheme")
	}
}

func TestOpenAPIJSONIsGeneratedFromServerRoutes(t *testing.T) {
	body, err := OpenAPIJSON()
	if err != nil {
		t.Fatalf("OpenAPIJSON error: %v", err)
	}
	if !strings.Contains(string(body), `"operationId": "websiteMetrics"`) {
		t.Fatalf("OpenAPI output does not include websiteMetrics operation")
	}
}
