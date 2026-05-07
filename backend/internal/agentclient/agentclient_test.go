package agentclient

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestClientAnalyzeSuccess(t *testing.T) {
	client := NewClient("http://agent.local")
	client.HTTP = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(`{"signal":"risk_warning","confidence":0.7,"risk_level":"medium","triggered_rules":["price_above"],"summary":"ok","reasoning":["test"],"data_time":"2026-05-05T00:00:00Z","source_refs":[],"missing_data":[],"recommended_action":"继续观察","indicators":{}}`), nil
	})}
	got, err := client.Analyze(context.Background(), AnalyzeRequest{Symbol: "AAPL", Market: "US", Prices: []float64{1, 2}})
	if err != nil {
		t.Fatalf("analyze: %v", err)
	}
	if got.Signal != "risk_warning" {
		t.Fatalf("unexpected result: %#v", got)
	}
}

func TestClientAnalyzeRejectsInvalidJSON(t *testing.T) {
	client := NewClient("http://agent.local")
	client.HTTP = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		return jsonResponse(`not-json`), nil
	})}
	if _, err := client.Analyze(context.Background(), AnalyzeRequest{}); err == nil {
		t.Fatal("expected invalid json error")
	}
}

func TestRunRepositorySaveAndFind(t *testing.T) {
	repo := NewRunRepository()
	run := Run{ID: "run-1", AgentVersion: "test", PromptVersion: "v1", Status: "succeeded"}

	if err := repo.Save(run); err != nil {
		t.Fatalf("save run: %v", err)
	}
	got, ok := repo.FindByID("run-1")
	if !ok || got.CreatedAt.IsZero() {
		t.Fatalf("run not saved correctly: %#v", got)
	}
}

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func jsonResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}
