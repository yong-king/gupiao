package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"jijin/backend/internal/config"
	"jijin/backend/internal/holdings"
	"jijin/backend/internal/watchlist"
)

func TestSystemDependenciesAPI(t *testing.T) {
	t.Setenv("DEEPSEEK_API_KEY", "test-key")
	cfg := config.Default()
	cfg.DatabaseURL = "postgres://jijin:jijin@127.0.0.1:1/jijin?sslmode=disable"
	cfg.RedisAddr = "127.0.0.1:1"
	server := NewServerWithConfig(watchlist.NewRepository(), holdings.NewRepository(), cfg)

	req := httptest.NewRequest(http.MethodGet, "/api/system/dependencies", nil)
	rec := httptest.NewRecorder()
	server.Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}
	var payload systemDependenciesResponse
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode dependencies: %v", err)
	}
	if payload.LLM.Name != "deepseek" || !payload.LLM.Reachable {
		t.Fatalf("unexpected llm status: %#v", payload.LLM)
	}
}
