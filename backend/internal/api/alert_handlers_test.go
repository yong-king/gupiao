package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"jijin/backend/internal/alerts"
	"jijin/backend/internal/holdings"
	"jijin/backend/internal/marketdata"
	"jijin/backend/internal/refresh"
	"jijin/backend/internal/watchlist"
)

func TestAlertRuleTriggersAfterManualRefresh(t *testing.T) {
	watchlists := watchlist.NewRepository()
	provider := marketdata.NewMockProvider()
	provider.SetQuote(marketdata.Snapshot{Market: "US", Symbol: "AAPL", Price: 150, Source: "mock", DataTime: time.Now().UTC()})
	service := refresh.NewService(provider, marketdata.NewSnapshotRepository(), refresh.NewJobRepository())
	server := NewServerWithRefresh(watchlists, holdings.NewRepository(), service)

	if err := watchlists.Save(watchlist.Watchlist{ID: "wl-1", UserID: "user-1", Name: "Core"}); err != nil {
		t.Fatalf("save watchlist: %v", err)
	}
	if err := watchlists.AddSymbol("wl-1", watchlist.Symbol{Market: "US", Symbol: "AAPL"}); err != nil {
		t.Fatalf("add symbol: %v", err)
	}

	ruleBody := `{"id":"rule-1","user_id":"user-1","market":"US","symbol":"AAPL","type":"price_below","threshold":160,"signal":"buy_watch","risk_level":"medium","enabled":true,"cooldown_seconds":1800}`
	ruleReq := httptest.NewRequest(http.MethodPost, "/api/alert-rules", bytes.NewBufferString(ruleBody))
	ruleRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(ruleRec, ruleReq)
	if ruleRec.Code != http.StatusCreated {
		t.Fatalf("expected rule status %d, got %d: %s", http.StatusCreated, ruleRec.Code, ruleRec.Body.String())
	}

	refreshReq := httptest.NewRequest(http.MethodPost, "/api/refresh/manual", bytes.NewBufferString(`{"user_id":"user-1","watchlist_id":"wl-1","job_id":"job-1"}`))
	refreshRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(refreshRec, refreshReq)
	if refreshRec.Code != http.StatusOK {
		t.Fatalf("expected refresh ok, got %d: %s", refreshRec.Code, refreshRec.Body.String())
	}

	alertReq := httptest.NewRequest(http.MethodGet, "/api/alerts?user_id=user-1", nil)
	alertRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(alertRec, alertReq)
	var events []alerts.Event
	if err := json.NewDecoder(alertRec.Body).Decode(&events); err != nil {
		t.Fatalf("decode alerts: %v", err)
	}
	if len(events) != 1 || events[0].Signal != alerts.SignalBuyWatch {
		t.Fatalf("unexpected events: %#v", events)
	}
}
