package watchlist

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"jijin/backend/internal/contracts"
)

type Symbol struct {
	Market    string  `json:"market"`
	Symbol    string  `json:"symbol"`
	Note      string  `json:"note"`
	BuyPrice  float64 `json:"buy_price"`
	SellPrice float64 `json:"sell_price"`
}

func NormalizeSymbol(symbol string) string {
	return strings.ToUpper(strings.TrimSpace(symbol))
}

func (s Symbol) Normalized() Symbol {
	s.Market = strings.ToUpper(strings.TrimSpace(s.Market))
	s.Symbol = NormalizeSymbol(s.Symbol)
	return s
}

func (s Symbol) Validate() error {
	normalized := s.Normalized()
	if normalized.Symbol == "" {
		return errors.New("symbol is required")
	}
	if !contracts.IsMarketCode(normalized.Market) {
		return fmt.Errorf("unsupported market %q", s.Market)
	}
	if s.BuyPrice < 0 || s.SellPrice < 0 {
		return errors.New("buy price and sell price must be non-negative")
	}
	return nil
}

func (s Symbol) Key() string {
	normalized := s.Normalized()
	return normalized.Market + ":" + normalized.Symbol
}

type Watchlist struct {
	ID        string
	UserID    string
	Name      string
	Symbols   []Symbol
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (w Watchlist) Validate() error {
	if w.ID == "" {
		return errors.New("watchlist id is required")
	}
	if w.UserID == "" {
		return errors.New("user id is required")
	}
	if strings.TrimSpace(w.Name) == "" {
		return errors.New("watchlist name is required")
	}
	return nil
}

type Repository struct {
	mu        sync.RWMutex
	watchlist map[string]Watchlist
}

func NewRepository() *Repository {
	return &Repository{watchlist: make(map[string]Watchlist)}
}

func (r *Repository) Save(w Watchlist) error {
	if err := w.Validate(); err != nil {
		return err
	}
	now := time.Now().UTC()
	if w.CreatedAt.IsZero() {
		w.CreatedAt = now
	}
	w.UpdatedAt = now

	r.mu.Lock()
	defer r.mu.Unlock()
	r.watchlist[w.ID] = copyWatchlist(w)
	return nil
}

func (r *Repository) FindByID(id string) (Watchlist, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	w, ok := r.watchlist[id]
	return copyWatchlist(w), ok
}

func (r *Repository) ListByUser(userID string) []Watchlist {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Watchlist, 0)
	for _, w := range r.watchlist {
		if w.UserID == userID {
			out = append(out, copyWatchlist(w))
		}
	}
	return out
}

func (r *Repository) AddSymbol(id string, symbol Symbol) error {
	if err := symbol.Validate(); err != nil {
		return err
	}
	symbol = symbol.Normalized()

	r.mu.Lock()
	defer r.mu.Unlock()

	w, ok := r.watchlist[id]
	if !ok {
		return errors.New("watchlist not found")
	}

	for _, existing := range w.Symbols {
		if existing.Key() == symbol.Key() {
			existing.Note = symbol.Note
			existing.BuyPrice = symbol.BuyPrice
			existing.SellPrice = symbol.SellPrice
			w.Symbols = replaceSymbol(w.Symbols, existing)
			w.UpdatedAt = time.Now().UTC()
			r.watchlist[id] = copyWatchlist(w)
			return nil
		}
	}

	w.Symbols = append(w.Symbols, symbol)
	w.UpdatedAt = time.Now().UTC()
	r.watchlist[id] = copyWatchlist(w)
	return nil
}

func replaceSymbol(symbols []Symbol, updated Symbol) []Symbol {
	out := append([]Symbol(nil), symbols...)
	for i, existing := range out {
		if existing.Key() == updated.Key() {
			out[i] = updated
			return out
		}
	}
	return append(out, updated)
}

func copyWatchlist(w Watchlist) Watchlist {
	w.Symbols = append([]Symbol(nil), w.Symbols...)
	return w
}
