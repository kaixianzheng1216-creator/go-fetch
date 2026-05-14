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
	if body := rec.Body.String(); !strings.Contains(body, `"detail":"未登录或登录已失效"`) {
		t.Fatalf("body = %q", body)
	}
}

func TestRoutesInvalidJSONReturnsProblemDetails(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader("{"))
	req.Header.Set("Content-Type", "application/json")
	testApp().Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", rec.Code)
	}
	if body := rec.Body.String(); !strings.Contains(body, `"status":400`) {
		t.Fatalf("body = %q", body)
	}
}

func TestRoutesValidateLoginBody(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	testApp().Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d", rec.Code)
	}
	if body := rec.Body.String(); !strings.Contains(body, `"detail":"validation failed"`) {
		t.Fatalf("body = %q", body)
	}
}

func TestRoutesValidateCollectBody(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/collect", strings.NewReader(`{"payload":{}}`))
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
		t.Fatalf("OpenAPI 内容缺少指标接口路径")
	}
	if !strings.Contains(body, `"sessionCookie"`) {
		t.Fatalf("OpenAPI 内容缺少 session cookie 安全方案")
	}
}

func TestOpenAPIJSONIsGeneratedFromServerRoutes(t *testing.T) {
	body, err := OpenAPIJSON()
	if err != nil {
		t.Fatalf("生成 OpenAPI JSON 失败: %v", err)
	}
	if !strings.Contains(string(body), `"operationId": "websiteMetrics"`) {
		t.Fatalf("OpenAPI 输出缺少 websiteMetrics 操作")
	}
}
