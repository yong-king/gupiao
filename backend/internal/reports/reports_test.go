package reports

import (
	"testing"
	"time"

	"jijin/backend/internal/alerts"
	"jijin/backend/internal/marketdata"
)

func TestBuildStockDetail(t *testing.T) {
	now := time.Now().UTC()
	price := float64(180)
	detail := BuildStockDetail("US", "AAPL",
		[]marketdata.Snapshot{{Market: "US", Symbol: "AAPL", Price: price, DataTime: now}},
		[]alerts.Event{{Market: "US", Symbol: "AAPL", Signal: alerts.SignalRiskWarning, RiskLevel: alerts.RiskHigh, CreatedAt: now}},
	)

	if detail.LatestPrice == nil || *detail.LatestPrice != 180 {
		t.Fatalf("unexpected latest price: %#v", detail.LatestPrice)
	}
	if detail.RiskLevel != alerts.RiskHigh {
		t.Fatalf("unexpected risk: %q", detail.RiskLevel)
	}
}

func TestGenerateDailyReport(t *testing.T) {
	now := time.Date(2026, 5, 5, 16, 0, 0, 0, time.UTC)
	report := GenerateDailyReport("report-1", "user-1", []alerts.Event{
		{UserID: "user-1", Market: "US", Symbol: "AAPL", Signal: alerts.SignalBuyWatch, Summary: "Price below threshold."},
	}, now)

	if len(report.RiskPoints) != 1 || len(report.NeedsConfirmation) != 1 {
		t.Fatalf("unexpected report: %#v", report)
	}
}

func TestGenerateDailyReportEmpty(t *testing.T) {
	report := GenerateDailyReport("report-1", "user-1", nil, time.Now().UTC())
	if report.Summary == "" || len(report.NeedsConfirmation) == 0 {
		t.Fatalf("unexpected empty report: %#v", report)
	}
}
