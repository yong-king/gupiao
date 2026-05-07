package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"jijin/backend/internal/holdings"
	"jijin/backend/internal/marketdata"
	"jijin/backend/internal/refresh"
	"jijin/backend/internal/watchlist"
)

func TestManualRefreshAPI(t *testing.T) {
	watchlists := watchlist.NewRepository()
	holdingRepo := holdings.NewRepository()
	provider := marketdata.NewMockProvider()
	provider.SetQuote(marketdata.Snapshot{Market: "US", Symbol: "AAPL", Price: 180, DataTime: time.Now().UTC()})
	service := refresh.NewService(provider, marketdata.NewSnapshotRepository(), refresh.NewJobRepository())
	server := NewServerWithRefresh(watchlists, holdingRepo, service)

	if err := watchlists.Save(watchlist.Watchlist{ID: "wl-1", UserID: "user-1", Name: "Core"}); err != nil {
		t.Fatalf("save watchlist: %v", err)
	}
	if err := watchlists.AddSymbol("wl-1", watchlist.Symbol{Market: "US", Symbol: "AAPL"}); err != nil {
		t.Fatalf("add symbol: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/refresh/manual", bytes.NewBufferString(`{
		"user_id":"user-1",
		"watchlist_id":"wl-1",
		"job_id":"job-1"
	}`))
	rec := httptest.NewRecorder()
	server.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}
	var job refresh.Job
	if err := json.NewDecoder(rec.Body).Decode(&job); err != nil {
		t.Fatalf("decode job: %v", err)
	}
	if job.Status != refresh.StatusSucceeded || len(job.Snapshots) != 1 {
		t.Fatalf("unexpected refresh job: %#v", job)
	}
}
