package holdings

import (
	"strings"
	"testing"
)

func TestParseCSVReturnsNormalizedHoldings(t *testing.T) {
	input := "market,symbol,quantity,cost_basis\nus, aapl ,10,150.50\nHK,0700,5,300\n"

	got, rowErrors, err := ParseCSV("user-1", strings.NewReader(input))
	if err != nil {
		t.Fatalf("parse csv: %v", err)
	}
	if len(rowErrors) != 0 {
		t.Fatalf("unexpected row errors: %#v", rowErrors)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 holdings, got %d", len(got))
	}
	if got[0].Market != "US" || got[0].Symbol != "AAPL" {
		t.Fatalf("holding was not normalized: %#v", got[0])
	}
	if got[0].AttentionLevel != "medium" {
		t.Fatalf("expected default attention level, got %#v", got[0])
	}
}

func TestParseCSVReportsRowErrors(t *testing.T) {
	input := "market,symbol,quantity,cost_basis\nUS,AAPL,nope,150.50\n"

	got, rowErrors, err := ParseCSV("user-1", strings.NewReader(input))
	if err != nil {
		t.Fatalf("parse csv: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected no valid holdings, got %d", len(got))
	}
	if len(rowErrors) != 1 {
		t.Fatalf("expected 1 row error, got %d", len(rowErrors))
	}
	if rowErrors[0].Row != 2 {
		t.Fatalf("expected row 2, got %d", rowErrors[0].Row)
	}
}

func TestRepositoryReplaceAndListByUser(t *testing.T) {
	repo := NewRepository()
	holdings := []Holding{
		{Market: "us", Symbol: "aapl", Quantity: 1, CostBasis: 100},
	}

	if err := repo.ReplaceForUser("user-1", holdings); err != nil {
		t.Fatalf("replace holdings: %v", err)
	}

	got := repo.ListByUser("user-1")
	if len(got) != 1 {
		t.Fatalf("expected 1 holding, got %d", len(got))
	}
	if got[0].UserID != "user-1" || got[0].Symbol != "AAPL" {
		t.Fatalf("unexpected holding: %#v", got[0])
	}
	if other := repo.ListByUser("user-2"); len(other) != 0 {
		t.Fatalf("expected no holdings for other user, got %d", len(other))
	}
}

func TestRepositoryUpsert(t *testing.T) {
	repo := NewRepository()
	if _, err := repo.Upsert(Holding{UserID: "user-1", Market: "us", Symbol: "aapl", Quantity: 10, CostBasis: 150}); err != nil {
		t.Fatalf("upsert holding: %v", err)
	}
	if _, err := repo.Upsert(Holding{UserID: "user-1", Market: "US", Symbol: "AAPL", Quantity: 12, CostBasis: 145}); err != nil {
		t.Fatalf("update holding: %v", err)
	}
	got := repo.ListByUser("user-1")
	if len(got) != 1 || got[0].Quantity != 12 || got[0].CostBasis != 145 {
		t.Fatalf("expected updated holding, got %#v", got)
	}
}

func TestRepositoryDelete(t *testing.T) {
	repo := NewRepository()
	if _, err := repo.Upsert(Holding{UserID: "user-1", Market: "CN", Symbol: "000821", Quantity: 100, CostBasis: 8}); err != nil {
		t.Fatalf("upsert holding: %v", err)
	}
	if !repo.Delete("user-1", "cn", "000821") {
		t.Fatal("expected delete to succeed")
	}
	if repo.Delete("user-1", "CN", "000821") {
		t.Fatal("expected second delete to report not found")
	}
	if got := repo.ListByUser("user-1"); len(got) != 0 {
		t.Fatalf("expected no holdings, got %#v", got)
	}
}
