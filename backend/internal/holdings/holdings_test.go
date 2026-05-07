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
