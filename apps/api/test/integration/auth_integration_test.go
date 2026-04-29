package integration

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"job-crawler/apps/api/internal/bootstrap"
)

func TestProtectedRouteRequiresToken(t *testing.T) {
	app := bootstrap.NewApp()
	req := httptest.NewRequest(http.MethodGet, "/api/feed", nil)
	res := httptest.NewRecorder()

	app.Router.ServeHTTP(res, req)
	if res.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without token, got %d", res.Code)
	}
}

func TestLoginEndpointReturnsSessionPayload(t *testing.T) {
	app := bootstrap.NewApp()
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"email":"mason@example.com","password":"password123"}`))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	app.Router.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Fatalf("expected 200 from login, got %d", res.Code)
	}
	if !strings.Contains(res.Body.String(), "\"token\"") {
		t.Fatalf("expected token in login response")
	}
}
