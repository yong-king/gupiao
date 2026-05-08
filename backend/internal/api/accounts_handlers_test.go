package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"jijin/backend/internal/accounts"
	"jijin/backend/internal/holdings"
	"jijin/backend/internal/watchlist"
)

func TestAccountsAPI(t *testing.T) {
	server := NewServer(watchlist.NewRepository(), holdings.NewRepository())
	body := `{"id":"acct-1","user_id":"user-1","alias":"只读账户","refresh_mode":"manual","read_only":true,"metadata":{"provider":"csv"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/accounts", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	server.Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected account create status %d, got %d: %s", http.StatusCreated, rec.Code, rec.Body.String())
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/accounts?user_id=user-1", nil)
	listRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(listRec, listReq)
	var got []accounts.Config
	if err := json.NewDecoder(listRec.Body).Decode(&got); err != nil {
		t.Fatalf("decode accounts: %v", err)
	}
	if len(got) != 1 || got[0].Alias != "只读账户" {
		t.Fatalf("unexpected accounts: %#v", got)
	}
}
