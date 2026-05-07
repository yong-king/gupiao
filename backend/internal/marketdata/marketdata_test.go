package marketdata

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestMockProviderFetchesQuotes(t *testing.T) {
	provider := NewMockProvider()
	provider.SetQuote(Snapshot{Market: "US", Symbol: "AAPL", Price: 180, DataTime: time.Now().UTC()})

	got, err := provider.FetchQuotes(context.Background(), []QuoteRequest{{Market: "US", Symbol: "AAPL"}})
	if err != nil {
		t.Fatalf("fetch quotes: %v", err)
	}
	if len(got) != 1 || got[0].Price != 180 {
		t.Fatalf("unexpected quotes: %#v", got)
	}
}

func TestMockProviderReturnsFailure(t *testing.T) {
	provider := NewMockProvider()
	provider.SetError(errors.New("provider down"))

	if _, err := provider.FetchQuotes(context.Background(), []QuoteRequest{{Market: "US", Symbol: "AAPL"}}); err == nil {
		t.Fatal("expected provider failure")
	}
}

func TestSnapshotRepositorySaveAndList(t *testing.T) {
	repo := NewSnapshotRepository()
	if err := repo.SaveAll([]Snapshot{{Market: "US", Symbol: "AAPL", Price: 180}}); err != nil {
		t.Fatalf("save snapshots: %v", err)
	}

	got := repo.ListBySymbol("US", "AAPL")
	if len(got) != 1 || got[0].Price != 180 {
		t.Fatalf("unexpected snapshots: %#v", got)
	}
}

func TestParseStooqCSV(t *testing.T) {
	csv := `Symbol,Name,Date,Time,Open,High,Low,Close,Volume
AAPL.US,APPLE INC,2026-05-06,15:53:31,281.915,285.39,281.07,285.37,3432889
`
	got, err := ParseStooqCSV(strings.NewReader(csv), QuoteRequest{Market: "US", Symbol: "AAPL"})
	if err != nil {
		t.Fatalf("parse stooq csv: %v", err)
	}
	if got.Source != "stooq" || got.Name != "APPLE INC" || got.Price != 285.37 || got.Open != 281.915 {
		t.Fatalf("unexpected snapshot: %#v", got)
	}
	if got.ChangePercent <= 0 {
		t.Fatalf("expected positive change percent, got %#v", got)
	}
}

func TestSnapshotRepositoryDailyChanges(t *testing.T) {
	repo := NewSnapshotRepository()
	if err := repo.SaveAll([]Snapshot{
		{Market: "US", Symbol: "AAPL", Open: 100, Price: 105, Volume: 10, Source: "test", DataTime: time.Date(2026, 5, 5, 15, 0, 0, 0, time.UTC)},
		{Market: "US", Symbol: "AAPL", Open: 106, Price: 103, Volume: 12, Source: "test", DataTime: time.Date(2026, 5, 6, 15, 0, 0, 0, time.UTC)},
	}); err != nil {
		t.Fatalf("save snapshots: %v", err)
	}
	changes := repo.DailyChanges("US", "AAPL")
	if len(changes) != 2 {
		t.Fatalf("expected 2 changes, got %#v", changes)
	}
	if changes[0].Change != 5 || changes[1].Change != -2 {
		t.Fatalf("unexpected changes: %#v", changes)
	}
	if !strings.Contains(changes[0].RAGText, "AAPL") {
		t.Fatalf("expected rag text, got %#v", changes[0])
	}
}

func TestProfileFromSnapshots(t *testing.T) {
	profile := ProfileFromSnapshots("US", "AAPL", []Snapshot{{Market: "US", Symbol: "AAPL", Name: "APPLE INC", ChangePercent: -4, DataTime: time.Now().UTC()}})
	if profile.Name != "APPLE INC" || profile.Sector == "" || profile.Recommendation != "risk_watch" {
		t.Fatalf("unexpected profile: %#v", profile)
	}
}
