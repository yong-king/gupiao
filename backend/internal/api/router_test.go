package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewRouterIncludesHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	NewRouter().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

func TestSwaggerOpenAPIIncludesAssistantAndWorkflowContracts(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/swagger/openapi.json", nil)
	rec := httptest.NewRecorder()

	NewRouter().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	var spec map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&spec); err != nil {
		t.Fatalf("decode swagger spec: %v", err)
	}
	paths, ok := spec["paths"].(map[string]interface{})
	if !ok {
		t.Fatalf("swagger paths missing: %#v", spec)
	}
	for _, path := range []string{"/api/workflows/research/run", "/api/assistant/chat", "/api/assistant/chat/stream", "/api/delivery/settings", "/api/delivery/test", "/api/delivery/logs"} {
		if _, ok := paths[path]; !ok {
			t.Fatalf("swagger path %s missing from %#v", path, paths)
		}
	}
	data, err := swaggerOpenAPIJSON()
	if err != nil || len(data) == 0 {
		t.Fatalf("swagger json helper failed: len=%d err=%v", len(data), err)
	}
}
