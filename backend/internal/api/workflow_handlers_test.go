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

func TestWorkflowRunProcessesAttentionLevelHoldings(t *testing.T) {
	holdingRepo := holdings.NewRepository()
	if _, err := holdingRepo.Upsert(holdings.Holding{UserID: "user-1", Market: "CN", Symbol: "000821", Quantity: 100, CostBasis: 8, AttentionLevel: "high"}); err != nil {
		t.Fatalf("upsert holding: %v", err)
	}
	provider := marketdata.NewMockProvider()
	provider.SetQuote(marketdata.Snapshot{Market: "CN", Symbol: "000821", Name: "京山轻机", Open: 10, Price: 10.5, PreviousClose: 10, ChangePercent: 5, Source: "mock", DataTime: time.Now().UTC()})
	service := refresh.NewService(provider, marketdata.NewSnapshotRepository(), refresh.NewJobRepository())
	server := NewServerWithRefresh(watchlist.NewRepository(), holdingRepo, service)

	req := httptest.NewRequest(http.MethodPost, "/api/workflows/research/run", strings.NewReader(`{"user_id":"user-1","attention_level":"high"}`))
	rec := httptest.NewRecorder()
	server.Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}
	var payload runWorkflowResponse
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload.Job.TargetCount != 1 || len(payload.Job.Steps) != 5 || len(payload.Documents) != 1 {
		t.Fatalf("unexpected workflow response: %#v", payload)
	}
}

func TestAssistantChatUsesHoldingContext(t *testing.T) {
	holdingRepo := holdings.NewRepository()
	if _, err := holdingRepo.Upsert(holdings.Holding{UserID: "user-1", Market: "CN", Symbol: "000821", Quantity: 100, CostBasis: 8, AttentionLevel: "high"}); err != nil {
		t.Fatalf("upsert holding: %v", err)
	}
	snapshots := marketdata.NewSnapshotRepository()
	if err := snapshots.SaveAll([]marketdata.Snapshot{{Market: "CN", Symbol: "000821", Name: "京山轻机", Price: 10.5, ChangePercent: 5, Source: "mock", DataTime: time.Now().UTC()}}); err != nil {
		t.Fatalf("save snapshot: %v", err)
	}
	service := refresh.NewService(marketdata.NewMockProvider(), snapshots, refresh.NewJobRepository())
	server := NewServerWithRefresh(watchlist.NewRepository(), holdingRepo, service)

	req := httptest.NewRequest(http.MethodPost, "/api/assistant/chat", strings.NewReader(`{"user_id":"user-1","symbol":"000821","question":"分析风险"}`))
	rec := httptest.NewRecorder()
	server.Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}
	var payload assistantChatResponse
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload.Market != "CN" || !strings.Contains(payload.ContextSummary, "持仓数量") || !strings.Contains(payload.Answer, "不是买卖指令") {
		t.Fatalf("unexpected assistant response: %#v", payload)
	}
}

func TestAssistantChatStreamEmitsChunks(t *testing.T) {
	holdingRepo := holdings.NewRepository()
	if _, err := holdingRepo.Upsert(holdings.Holding{UserID: "user-1", Market: "CN", Symbol: "000821", Quantity: 100, CostBasis: 8, AttentionLevel: "high"}); err != nil {
		t.Fatalf("upsert holding: %v", err)
	}
	server := NewServerWithRefresh(watchlist.NewRepository(), holdingRepo, refresh.NewService(marketdata.NewMockProvider(), marketdata.NewSnapshotRepository(), refresh.NewJobRepository()))

	req := httptest.NewRequest(http.MethodPost, "/api/assistant/chat/stream", strings.NewReader(`{"user_id":"user-1","session_id":"s1","symbol":"000821","question":"分析风险"}`))
	rec := httptest.NewRecorder()
	server.Routes().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d: %s", http.StatusOK, rec.Code, rec.Body.String())
	}
	if contentType := rec.Header().Get("Content-Type"); !strings.Contains(contentType, "text/event-stream") {
		t.Fatalf("expected event stream, got %q", contentType)
	}
	if body := rec.Body.String(); !strings.Contains(body, `"delta"`) || !strings.Contains(body, `"done":true`) {
		t.Fatalf("unexpected stream body: %s", body)
	}
}
