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

func TestWatchlistAPIFlow(t *testing.T) {
	server := NewServer(watchlist.NewRepository(), holdings.NewRepository())

	createReq := httptest.NewRequest(http.MethodPost, "/api/watchlists", bytes.NewBufferString(`{
		"id":"wl-1",
		"user_id":"user-1",
		"name":"Core"
	}`))
	createRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(createRec, createReq)
	if createRec.Code != http.StatusCreated {
		t.Fatalf("expected create status %d, got %d: %s", http.StatusCreated, createRec.Code, createRec.Body.String())
	}

	addReq := httptest.NewRequest(http.MethodPost, "/api/watchlists/wl-1/symbols", bytes.NewBufferString(`{
		"market":"us",
		"symbol":" aapl ",
		"buy_price":150,
		"sell_price":220
	}`))
	addRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(addRec, addReq)
	if addRec.Code != http.StatusOK {
		t.Fatalf("expected add status %d, got %d: %s", http.StatusOK, addRec.Code, addRec.Body.String())
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/watchlists/wl-1", nil)
	getRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(getRec, getReq)
	if getRec.Code != http.StatusOK {
		t.Fatalf("expected get status %d, got %d", http.StatusOK, getRec.Code)
	}

	var got watchlist.Watchlist
	if err := json.NewDecoder(getRec.Body).Decode(&got); err != nil {
		t.Fatalf("decode watchlist: %v", err)
	}
	if len(got.Symbols) != 1 || got.Symbols[0].Symbol != "AAPL" || got.Symbols[0].BuyPrice != 150 {
		t.Fatalf("unexpected symbols: %#v", got.Symbols)
	}
	if audits := server.AuditEntries(); len(audits) != 2 {
		t.Fatalf("expected 2 audit entries, got %d", len(audits))
	}
}

func TestHoldingsImportAPI(t *testing.T) {
	server := NewServer(watchlist.NewRepository(), holdings.NewRepository())
	body := `{
		"user_id":"user-1",
		"csv":"market,symbol,quantity,cost_basis\nUS,AAPL,10,150\nUS,MSFT,nope,300\n"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/holdings/import", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	server.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected import status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	var got importHoldingsResponse
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Imported != 1 || len(got.RowErrors) != 1 {
		t.Fatalf("unexpected import response: %#v", got)
	}
	if audits := server.AuditEntries(); len(audits) != 1 || audits[0].Action != "holdings.import" {
		t.Fatalf("unexpected audit entries: %#v", audits)
	}
}

func TestHoldingsUpsertAPI(t *testing.T) {
	server := NewServer(watchlist.NewRepository(), holdings.NewRepository())
	body := `{"user_id":"user-1","market":"US","symbol":"AAPL","quantity":10,"cost_basis":150}`

	req := httptest.NewRequest(http.MethodPost, "/api/holdings", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()
	server.Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected upsert status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}

	listReq := httptest.NewRequest(http.MethodGet, "/api/holdings?user_id=user-1", nil)
	listRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(listRec, listReq)
	if listRec.Code != http.StatusOK {
		t.Fatalf("expected list status %d, got %d", http.StatusOK, listRec.Code)
	}
	var got []holdings.Holding
	if err := json.NewDecoder(listRec.Body).Decode(&got); err != nil {
		t.Fatalf("decode holdings: %v", err)
	}
	if len(got) != 1 || got[0].Symbol != "AAPL" || got[0].Quantity != 10 {
		t.Fatalf("unexpected holdings: %#v", got)
	}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/holdings?user_id=user-1&market=US&symbol=AAPL", nil)
	deleteRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(deleteRec, deleteReq)
	if deleteRec.Code != http.StatusOK {
		t.Fatalf("expected delete status %d, got %d: %s", http.StatusOK, deleteRec.Code, deleteRec.Body.String())
	}

	emptyReq := httptest.NewRequest(http.MethodGet, "/api/holdings?user_id=user-1", nil)
	emptyRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(emptyRec, emptyReq)
	if emptyRec.Code != http.StatusOK {
		t.Fatalf("expected list status %d, got %d", http.StatusOK, emptyRec.Code)
	}
	var empty []holdings.Holding
	if err := json.NewDecoder(emptyRec.Body).Decode(&empty); err != nil {
		t.Fatalf("decode holdings: %v", err)
	}
	if len(empty) != 0 {
		t.Fatalf("expected holdings to be deleted, got %#v", empty)
	}
}
