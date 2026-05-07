package watchlist

import "testing"

func TestRepositorySavesAndFindsWatchlist(t *testing.T) {
	repo := NewRepository()
	want := Watchlist{ID: "wl-1", UserID: "user-1", Name: "Core"}

	if err := repo.Save(want); err != nil {
		t.Fatalf("save watchlist: %v", err)
	}

	got, ok := repo.FindByID("wl-1")
	if !ok {
		t.Fatal("expected watchlist to be found")
	}
	if got.ID != want.ID || got.UserID != want.UserID || got.Name != want.Name {
		t.Fatalf("unexpected watchlist: %#v", got)
	}
}

func TestAddSymbolNormalizesAndUpdatesDuplicate(t *testing.T) {
	repo := NewRepository()
	if err := repo.Save(Watchlist{ID: "wl-1", UserID: "user-1", Name: "Core"}); err != nil {
		t.Fatalf("save watchlist: %v", err)
	}

	if err := repo.AddSymbol("wl-1", Symbol{Market: "us", Symbol: " aapl ", BuyPrice: 150, SellPrice: 220}); err != nil {
		t.Fatalf("add symbol: %v", err)
	}
	if err := repo.AddSymbol("wl-1", Symbol{Market: "US", Symbol: "AAPL", BuyPrice: 160, SellPrice: 230}); err != nil {
		t.Fatalf("update duplicate symbol: %v", err)
	}

	got, _ := repo.FindByID("wl-1")
	if len(got.Symbols) != 1 || got.Symbols[0].Market != "US" || got.Symbols[0].Symbol != "AAPL" {
		t.Fatalf("symbol was not normalized: %#v", got.Symbols[0])
	}
	if got.Symbols[0].BuyPrice != 160 || got.Symbols[0].SellPrice != 230 {
		t.Fatalf("monitor prices were not updated: %#v", got.Symbols[0])
	}
}

func TestSymbolRejectsUnsupportedMarket(t *testing.T) {
	if err := (Symbol{Market: "AUTO", Symbol: "AAPL"}).Validate(); err == nil {
		t.Fatal("expected unsupported market to fail")
	}
}

func TestSymbolRejectsNegativeMonitorPrices(t *testing.T) {
	if err := (Symbol{Market: "US", Symbol: "AAPL", BuyPrice: -1}).Validate(); err == nil {
		t.Fatal("expected negative monitor price to fail")
	}
}
