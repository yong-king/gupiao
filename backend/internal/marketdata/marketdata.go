package marketdata

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"jijin/backend/internal/watchlist"
)

type QuoteRequest struct {
	Market string
	Symbol string
}

func FromWatchlistSymbol(symbol watchlist.Symbol) QuoteRequest {
	normalized := symbol.Normalized()
	return QuoteRequest{Market: normalized.Market, Symbol: normalized.Symbol}
}

func (r QuoteRequest) Key() string {
	return r.Market + ":" + r.Symbol
}

type Snapshot struct {
	Market        string
	Symbol        string
	Name          string
	Open          float64
	High          float64
	Low           float64
	Price         float64
	PreviousClose float64
	ChangePercent float64
	Volume        int64
	Source        string
	DataTime      time.Time
	CreatedAt     time.Time
}

type Provider interface {
	FetchQuotes(ctx context.Context, requests []QuoteRequest) ([]Snapshot, error)
}

type MockProvider struct {
	mu     sync.RWMutex
	quotes map[string]Snapshot
	err    error
}

func NewMockProvider() *MockProvider {
	return &MockProvider{quotes: make(map[string]Snapshot)}
}

func (p *MockProvider) SetQuote(snapshot Snapshot) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if snapshot.Source == "" {
		snapshot.Source = "mock"
	}
	if snapshot.DataTime.IsZero() {
		snapshot.DataTime = time.Now().UTC()
	}
	p.quotes[QuoteRequest{Market: snapshot.Market, Symbol: snapshot.Symbol}.Key()] = snapshot
}

func (p *MockProvider) SetError(err error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.err = err
}

func (p *MockProvider) FetchQuotes(ctx context.Context, requests []QuoteRequest) ([]Snapshot, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if p.err != nil {
		return nil, p.err
	}

	snapshots := make([]Snapshot, 0, len(requests))
	for _, request := range requests {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		snapshot, ok := p.quotes[request.Key()]
		if !ok {
			return nil, fmt.Errorf("quote not found for %s", request.Key())
		}
		snapshot.CreatedAt = time.Now().UTC()
		snapshots = append(snapshots, snapshot)
	}
	return snapshots, nil
}

type SnapshotRepository struct {
	mu        sync.RWMutex
	snapshots []Snapshot
}

func NewSnapshotRepository() *SnapshotRepository {
	return &SnapshotRepository{}
}

func (r *SnapshotRepository) SaveAll(snapshots []Snapshot) error {
	for _, snapshot := range snapshots {
		if snapshot.Market == "" || snapshot.Symbol == "" {
			return errors.New("snapshot market and symbol are required")
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.snapshots = append(r.snapshots, append([]Snapshot(nil), snapshots...)...)
	return nil
}

func (r *SnapshotRepository) ListBySymbol(market string, symbol string) []Snapshot {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []Snapshot
	for _, snapshot := range r.snapshots {
		if snapshot.Market == market && snapshot.Symbol == symbol {
			out = append(out, snapshot)
		}
	}
	return out
}

type DailyChange struct {
	Market        string
	Symbol        string
	Date          string
	Open          float64
	Close         float64
	PreviousClose float64
	Change        float64
	ChangePercent float64
	Volume        int64
	Source        string
	DataTime      time.Time
	RAGText       string
}

func (r *SnapshotRepository) DailyChanges(market string, symbol string) []DailyChange {
	snapshots := r.ListBySymbol(market, symbol)
	byDate := make(map[string]Snapshot)
	var dates []string
	for _, snapshot := range snapshots {
		date := snapshot.DataTime.Format("2006-01-02")
		if date == "0001-01-01" {
			date = snapshot.CreatedAt.Format("2006-01-02")
		}
		if existing, ok := byDate[date]; !ok {
			dates = append(dates, date)
			byDate[date] = snapshot
		} else if snapshot.DataTime.After(existing.DataTime) {
			byDate[date] = snapshot
		}
	}
	sortStrings(dates)

	changes := make([]DailyChange, 0, len(dates))
	var previousClose float64
	for _, date := range dates {
		snapshot := byDate[date]
		baseline := snapshot.PreviousClose
		if baseline == 0 {
			baseline = previousClose
		}
		if baseline == 0 {
			baseline = snapshot.Open
		}
		change := snapshot.Price - baseline
		changePercent := 0.0
		if baseline != 0 {
			changePercent = change / baseline * 100
		}
		record := DailyChange{
			Market:        snapshot.Market,
			Symbol:        snapshot.Symbol,
			Date:          date,
			Open:          snapshot.Open,
			Close:         snapshot.Price,
			PreviousClose: baseline,
			Change:        change,
			ChangePercent: changePercent,
			Volume:        snapshot.Volume,
			Source:        snapshot.Source,
			DataTime:      snapshot.DataTime,
		}
		record.RAGText = fmt.Sprintf("%s %s %s close %.2f, change %.2f, change_percent %.2f%%, volume %d, source %s.",
			record.Date, record.Market, record.Symbol, record.Close, record.Change, record.ChangePercent, record.Volume, record.Source)
		changes = append(changes, record)
		previousClose = snapshot.Price
	}
	return changes
}

func (r *SnapshotRepository) ListAll() []Snapshot {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]Snapshot(nil), r.snapshots...)
}

func sortStrings(values []string) {
	for i := 1; i < len(values); i++ {
		value := values[i]
		j := i - 1
		for ; j >= 0 && values[j] > value; j-- {
			values[j+1] = values[j]
		}
		values[j+1] = value
	}
}
