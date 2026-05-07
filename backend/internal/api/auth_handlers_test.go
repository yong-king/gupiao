package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"jijin/backend/internal/holdings"
	"jijin/backend/internal/watchlist"
)

func TestRegisterLoginAndProtectedAPI(t *testing.T) {
	server := NewServer(watchlist.NewRepository(), holdings.NewRepository())
	server.authRequired = true

	registerReq := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBufferString(`{"id":"user-1","email":"test@example.com","password":"password123"}`))
	registerRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(registerRec, registerReq)
	if registerRec.Code != http.StatusCreated {
		t.Fatalf("expected register status %d, got %d: %s", http.StatusCreated, registerRec.Code, registerRec.Body.String())
	}
	var registered authResponse
	if err := json.NewDecoder(registerRec.Body).Decode(&registered); err != nil {
		t.Fatalf("decode register: %v", err)
	}

	blockedReq := httptest.NewRequest(http.MethodPost, "/api/watchlists", bytes.NewBufferString(`{"id":"wl-1","user_id":"user-1","name":"Core"}`))
	blockedRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(blockedRec, blockedReq)
	if blockedRec.Code != http.StatusUnauthorized {
		t.Fatalf("expected unauthorized, got %d", blockedRec.Code)
	}

	createReq := httptest.NewRequest(http.MethodPost, "/api/watchlists", bytes.NewBufferString(`{"id":"wl-1","user_id":"user-1","name":"Core"}`))
	createReq.Header.Set("Authorization", "Bearer "+registered.Token)
	createRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(createRec, createReq)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected created with token, got %d: %s", createRec.Code, createRec.Body.String())
	}
}
