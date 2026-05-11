package api

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"jijin/backend/internal/alerts"
	"jijin/backend/internal/delivery"
	"jijin/backend/internal/holdings"
	"jijin/backend/internal/marketdata"
	"jijin/backend/internal/refresh"
	"jijin/backend/internal/watchlist"
)

type testSender struct {
	response string
	err      error
}

func (s *testSender) Send(ctx context.Context, settings delivery.Settings, title string, body string) (string, error) {
	return s.response, s.err
}

func TestDeliverySettingsAndTestMessage(t *testing.T) {
	server := NewServer(watchlist.NewRepository(), holdings.NewRepository())
	server.wecom = &testSender{response: "ok"}

	saveReq := httptest.NewRequest(http.MethodPost, "/api/delivery/settings", bytes.NewBufferString(`{
		"user_id":"user-1",
		"enabled":true,
		"default_channel":"wecom_webhook",
		"wecom_webhook_url":"https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=test",
		"wecom_mention":"所有人",
		"min_risk_level":"high",
		"cooldown_seconds":1800,
		"mcp_enabled":false
	}`))
	saveRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(saveRec, saveReq)
	if saveRec.Code != http.StatusOK {
		t.Fatalf("expected settings save ok, got %d: %s", saveRec.Code, saveRec.Body.String())
	}

	testReq := httptest.NewRequest(http.MethodPost, "/api/delivery/test", bytes.NewBufferString(`{"user_id":"user-1"}`))
	testRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(testRec, testReq)
	if testRec.Code != http.StatusOK {
		t.Fatalf("expected delivery test ok, got %d: %s", testRec.Code, testRec.Body.String())
	}

	logReq := httptest.NewRequest(http.MethodGet, "/api/delivery/logs?user_id=user-1", nil)
	logRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(logRec, logReq)
	if logRec.Code != http.StatusOK {
		t.Fatalf("expected delivery logs ok, got %d: %s", logRec.Code, logRec.Body.String())
	}
}

func TestHighRiskAlertCreatesDeliveryLog(t *testing.T) {
	watchlists := watchlist.NewRepository()
	provider := marketdata.NewMockProvider()
	provider.SetQuote(marketdata.Snapshot{Market: "US", Symbol: "AAPL", Price: 150, Source: "mock", DataTime: time.Now().UTC()})
	service := refresh.NewService(provider, marketdata.NewSnapshotRepository(), refresh.NewJobRepository())
	server := NewServerWithRefresh(watchlists, holdings.NewRepository(), service)
	server.wecom = &testSender{response: "delivered"}

	if err := watchlists.Save(watchlist.Watchlist{ID: "wl-high", UserID: "user-1", Name: "Core"}); err != nil {
		t.Fatalf("save watchlist: %v", err)
	}
	if err := watchlists.AddSymbol("wl-high", watchlist.Symbol{Market: "US", Symbol: "AAPL"}); err != nil {
		t.Fatalf("add symbol: %v", err)
	}
	if err := server.deliveries.Save(delivery.Settings{
		UserID:          "user-1",
		Enabled:         true,
		DefaultChannel:  delivery.ChannelTypeWeCom,
		WeComWebhookURL: "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=test",
		MinRiskLevel:    string(alerts.RiskHigh),
		CooldownSeconds: 1800,
	}); err != nil {
		t.Fatalf("save delivery settings: %v", err)
	}

	ruleReq := httptest.NewRequest(http.MethodPost, "/api/alert-rules", bytes.NewBufferString(`{
		"id":"rule-high",
		"user_id":"user-1",
		"market":"US",
		"symbol":"AAPL",
		"type":"price_below",
		"threshold":160,
		"signal":"risk_warning",
		"risk_level":"high",
		"enabled":true,
		"cooldown_seconds":1800
	}`))
	ruleRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(ruleRec, ruleReq)
	if ruleRec.Code != http.StatusCreated {
		t.Fatalf("expected rule created, got %d: %s", ruleRec.Code, ruleRec.Body.String())
	}

	refreshReq := httptest.NewRequest(http.MethodPost, "/api/refresh/manual", bytes.NewBufferString(`{"user_id":"user-1","watchlist_id":"wl-high","job_id":"job-high"}`))
	refreshRec := httptest.NewRecorder()
	server.Routes().ServeHTTP(refreshRec, refreshReq)
	if refreshRec.Code != http.StatusOK {
		t.Fatalf("expected refresh ok, got %d: %s", refreshRec.Code, refreshRec.Body.String())
	}

	logs := server.deliveryLogs.ListByUser("user-1", 10)
	if len(logs) == 0 || logs[0].Status != "sent" {
		t.Fatalf("expected delivery log, got %#v", logs)
	}
}
