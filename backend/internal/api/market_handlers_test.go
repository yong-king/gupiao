package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"jijin/backend/internal/holdings"
	"jijin/backend/internal/marketdata"
	"jijin/backend/internal/refresh"
	"jijin/backend/internal/watchlist"
)

func TestMarketSnapshotsDailyChangesAndProfileAPI(t *testing.T) {
	snapshots := marketdata.NewSnapshotRepository()
	if err := snapshots.SaveAll([]marketdata.Snapshot{{
		Market: "US", Symbol: "AAPL", Name: "APPLE INC", Open: 100, Price: 105, Source: "test", DataTime: time.Date(2026, 5, 6, 15, 0, 0, 0, time.UTC),
	}}); err != nil {
		t.Fatalf("save snapshot: %v", err)
	}
	service := refresh.NewService(marketdata.NewMockProvider(), snapshots, refresh.NewJobRepository())
	server := NewServerWithRefresh(watchlist.NewRepository(), holdings.NewRepository(), service)

	for _, path := range []string{
		"/api/market/snapshots?market=US&symbol=AAPL",
		"/api/market/daily-changes?market=US&symbol=AAPL",
		"/api/stocks/profile?market=US&symbol=AAPL",
	} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		server.Routes().ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Fatalf("%s expected status %d, got %d: %s", path, http.StatusOK, rec.Code, rec.Body.String())
		}
		var payload any
		if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
			t.Fatalf("decode %s: %v", path, err)
		}
	}
}

func TestMarketCollectAPI(t *testing.T) {
	provider := marketdata.NewMockProvider()
	provider.SetQuote(marketdata.Snapshot{
		Market: "US", Symbol: "AAPL", Name: "APPLE INC", Open: 100, Price: 106, Source: "mock", DataTime: time.Date(2026, 5, 6, 15, 0, 0, 0, time.UTC),
	})
	service := refresh.NewService(provider, marketdata.NewSnapshotRepository(), refresh.NewJobRepository())
	server := NewServerWithRefresh(watchlist.NewRepository(), holdings.NewRepository(), service)

	req := httptest.NewRequest(http.MethodPost, "/api/market/collect", strings.NewReader(`{"market":"US","symbol":"AAPL"}`))
	rec := httptest.NewRecorder()
	server.Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}
	var payload collectMarketResponse
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode collect response: %v", err)
	}
	if payload.Snapshot.Price != 106 || len(payload.DailyChanges) != 1 || payload.Profile.Name != "APPLE INC" {
		t.Fatalf("unexpected payload: %#v", payload)
	}
}

func TestMarketCollectFallsBackToSavedSnapshot(t *testing.T) {
	provider := marketdata.NewMockProvider()
	provider.SetError(assertiveError("eastmoney EOF"))
	snapshots := marketdata.NewSnapshotRepository()
	if err := snapshots.SaveAll([]marketdata.Snapshot{{
		Market: "CN", Symbol: "000821", Name: "京山轻机", Price: 10.71, PreviousClose: 10.50, ChangePercent: 2, Source: "saved", DataTime: time.Date(2026, 5, 9, 9, 30, 0, 0, time.UTC),
	}}); err != nil {
		t.Fatalf("save snapshot: %v", err)
	}
	service := refresh.NewService(provider, snapshots, refresh.NewJobRepository())
	server := NewServerWithRefresh(watchlist.NewRepository(), holdings.NewRepository(), service)

	req := httptest.NewRequest(http.MethodPost, "/api/market/collect", strings.NewReader(`{"market":"CN","symbol":"000821"}`))
	rec := httptest.NewRecorder()
	server.Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected fallback status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}
	var payload collectMarketResponse
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode collect response: %v", err)
	}
	if payload.Snapshot.Price != 10.71 || payload.Warning == "" {
		t.Fatalf("expected saved quote fallback, got %#v", payload)
	}
}

func TestMarketKLinesFallsBackToSnapshots(t *testing.T) {
	snapshots := marketdata.NewSnapshotRepository()
	if err := snapshots.SaveAll([]marketdata.Snapshot{{
		Market: "US", Symbol: "AAPL", Open: 100, High: 106, Low: 99, Price: 105, Volume: 10, Source: "saved", DataTime: time.Date(2026, 5, 8, 0, 0, 0, 0, time.UTC),
	}}); err != nil {
		t.Fatalf("save snapshot: %v", err)
	}
	service := refresh.NewService(marketdata.NewMockProvider(), snapshots, refresh.NewJobRepository())
	server := NewServerWithRefresh(watchlist.NewRepository(), holdings.NewRepository(), service)

	req := httptest.NewRequest(http.MethodGet, "/api/market/klines?market=US&symbol=AAPL", nil)
	rec := httptest.NewRecorder()
	server.Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}
	var payload []marketdata.KLine
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode klines: %v", err)
	}
	if len(payload) != 1 || payload[0].Close != 105 || payload[0].High != 106 {
		t.Fatalf("unexpected klines: %#v", payload)
	}
}

func TestResearchCollectAPI(t *testing.T) {
	snapshots := marketdata.NewSnapshotRepository()
	if err := snapshots.SaveAll([]marketdata.Snapshot{{
		Market: "CN", Symbol: "000821", Name: "京山轻机", Open: 10, Price: 10.5, PreviousClose: 10, ChangePercent: 5, Source: "test", DataTime: time.Date(2026, 5, 8, 15, 0, 0, 0, time.UTC),
	}}); err != nil {
		t.Fatalf("save snapshot: %v", err)
	}
	service := refresh.NewService(marketdata.NewMockProvider(), snapshots, refresh.NewJobRepository())
	server := NewServerWithRefresh(watchlist.NewRepository(), holdings.NewRepository(), service)

	req := httptest.NewRequest(http.MethodPost, "/api/research/collect", strings.NewReader(`{"user_id":"user-1","market":"CN","symbol":"000821","attention_level":"high"}`))
	rec := httptest.NewRecorder()
	server.Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}
	var payload collectResearchResponse
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode research response: %v", err)
	}
	if payload.AttentionLevel != "high" || !strings.Contains(payload.Summary, "1h0m0s") {
		t.Fatalf("unexpected payload: %#v", payload)
	}
}

type assertiveError string

func (e assertiveError) Error() string {
	return string(e)
}
