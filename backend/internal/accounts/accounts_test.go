package accounts

import (
	"testing"

	"jijin/backend/internal/holdings"
	"jijin/backend/internal/marketdata"
	"jijin/backend/internal/settings"
)

func TestConfigValidationRejectsSensitiveMetadata(t *testing.T) {
	config := Config{ID: "acct-1", UserID: "user-1", Alias: "Main", RefreshMode: settings.RefreshModeConservative, ReadOnly: true, Metadata: map[string]string{"password": "secret"}}
	if err := config.Validate(); err == nil {
		t.Fatal("expected sensitive metadata to fail")
	}
}

func TestConfigMustBeReadOnly(t *testing.T) {
	config := Config{ID: "acct-1", UserID: "user-1", Alias: "Main", RefreshMode: settings.RefreshModeConservative, ReadOnly: false}
	if err := config.Validate(); err == nil {
		t.Fatal("expected non-read-only account to fail")
	}
}

func TestCalculateRisk(t *testing.T) {
	result := CalculateRisk(
		[]holdings.Holding{
			{Market: "US", Symbol: "AAPL", Quantity: 10, CostBasis: 100},
			{Market: "US", Symbol: "MSFT", Quantity: 1, CostBasis: 100},
		},
		[]marketdata.Snapshot{
			{Market: "US", Symbol: "AAPL", Price: 200},
			{Market: "US", Symbol: "MSFT", Price: 100},
		},
		0.8,
	)

	if result.TotalValue != 2100 {
		t.Fatalf("unexpected total value %f", result.TotalValue)
	}
	if result.Positions[0].RiskObservation == "继续观察" {
		t.Fatalf("expected concentration risk observation: %#v", result.Positions[0])
	}
}

func TestCalculateRiskMissingPrice(t *testing.T) {
	result := CalculateRisk([]holdings.Holding{{Market: "US", Symbol: "AAPL", Quantity: 1, CostBasis: 100}}, nil, 0.8)
	if len(result.MissingData) != 1 {
		t.Fatalf("expected missing data: %#v", result)
	}
}
