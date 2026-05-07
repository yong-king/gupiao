package stocksource

import (
	"testing"
)

func TestParseJSON(t *testing.T) {
	got, err := ParseJSON([]byte(`{"market":"US","symbol":"AAPL","price":180,"previous_close":175,"change_percent":2.8,"volume":1000,"source":"demo","data_time":"2026-05-05T00:00:00Z"}`))
	if err != nil {
		t.Fatalf("parse json: %v", err)
	}
	if got.Symbol != "AAPL" || got.Price != 180 || got.Source != "demo" {
		t.Fatalf("unexpected snapshot: %#v", got)
	}
}

func TestParseHTML(t *testing.T) {
	html := `<div data-market="US" data-symbol="AAPL" data-price="180" data-previous-close="175" data-change-percent="2.8" data-volume="1000" data-source="demo" data-data-time="2026-05-05T00:00:00Z"></div>`
	got, err := ParseHTML(html)
	if err != nil {
		t.Fatalf("parse html: %v", err)
	}
	if got.Symbol != "AAPL" || got.Price != 180 {
		t.Fatalf("unexpected snapshot: %#v", got)
	}
}

func TestParseJSONMissingFields(t *testing.T) {
	if _, err := ParseJSON([]byte(`{"market":"US"}`)); err == nil {
		t.Fatal("expected missing fields to fail")
	}
}
