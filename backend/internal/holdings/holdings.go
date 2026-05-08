package holdings

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"jijin/backend/internal/contracts"
	"jijin/backend/internal/watchlist"
)

type Holding struct {
	ID             string
	UserID         string
	Market         string
	Symbol         string
	Quantity       float64
	CostBasis      float64
	AttentionLevel string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (h Holding) Normalized() Holding {
	h.Market = strings.ToUpper(strings.TrimSpace(h.Market))
	h.Symbol = watchlist.NormalizeSymbol(h.Symbol)
	h.AttentionLevel = NormalizeAttentionLevel(h.AttentionLevel)
	return h
}

func (h Holding) Validate() error {
	normalized := h.Normalized()
	if normalized.UserID == "" {
		return errors.New("user id is required")
	}
	if normalized.Symbol == "" {
		return errors.New("symbol is required")
	}
	if !contracts.IsMarketCode(normalized.Market) {
		return fmt.Errorf("unsupported market %q", h.Market)
	}
	if normalized.Quantity <= 0 {
		return errors.New("quantity must be greater than zero")
	}
	if normalized.CostBasis < 0 {
		return errors.New("cost basis must not be negative")
	}
	switch normalized.AttentionLevel {
	case "low", "medium", "high":
	default:
		return errors.New("attention level must be low, medium or high")
	}
	return nil
}

func NormalizeAttentionLevel(level string) string {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "high":
		return "high"
	case "low":
		return "low"
	default:
		return "medium"
	}
}

func AttentionRefreshInterval(level string) time.Duration {
	switch NormalizeAttentionLevel(level) {
	case "high":
		return 4 * time.Hour
	case "low":
		return 24 * time.Hour
	default:
		return 6 * time.Hour
	}
}

type RowError struct {
	Row int
	Err error
}

func (e RowError) Error() string {
	return fmt.Sprintf("row %d: %v", e.Row, e.Err)
}

func ParseCSV(userID string, reader io.Reader) ([]Holding, []RowError, error) {
	csvReader := csv.NewReader(reader)
	csvReader.TrimLeadingSpace = true

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, nil, err
	}
	if len(records) == 0 {
		return nil, nil, errors.New("csv is empty")
	}

	headers := normalizeHeaders(records[0])
	index := map[string]int{}
	for i, header := range headers {
		index[header] = i
	}
	for _, required := range []string{"market", "symbol", "quantity", "cost_basis"} {
		if _, ok := index[required]; !ok {
			return nil, nil, fmt.Errorf("missing required header %q", required)
		}
	}

	var holdings []Holding
	var rowErrors []RowError
	for i, record := range records[1:] {
		rowNumber := i + 2
		holding, err := parseRecord(userID, record, index)
		if err != nil {
			rowErrors = append(rowErrors, RowError{Row: rowNumber, Err: err})
			continue
		}
		holdings = append(holdings, holding.Normalized())
	}

	return holdings, rowErrors, nil
}

func parseRecord(userID string, record []string, index map[string]int) (Holding, error) {
	value := func(name string) string {
		i, ok := index[name]
		if !ok {
			return ""
		}
		if i >= len(record) {
			return ""
		}
		return strings.TrimSpace(record[i])
	}

	quantity, err := strconv.ParseFloat(value("quantity"), 64)
	if err != nil {
		return Holding{}, fmt.Errorf("invalid quantity: %w", err)
	}
	costBasis, err := strconv.ParseFloat(value("cost_basis"), 64)
	if err != nil {
		return Holding{}, fmt.Errorf("invalid cost_basis: %w", err)
	}

	holding := Holding{
		UserID:         userID,
		Market:         value("market"),
		Symbol:         value("symbol"),
		Quantity:       quantity,
		CostBasis:      costBasis,
		AttentionLevel: value("attention_level"),
	}
	if err := holding.Validate(); err != nil {
		return Holding{}, err
	}
	return holding, nil
}

func normalizeHeaders(headers []string) []string {
	out := make([]string, len(headers))
	for i, header := range headers {
		out[i] = strings.ToLower(strings.TrimSpace(header))
	}
	return out
}

type Repository struct {
	mu       sync.RWMutex
	holdings map[string][]Holding
}

func NewRepository() *Repository {
	return &Repository{holdings: make(map[string][]Holding)}
}

func (r *Repository) ReplaceForUser(userID string, holdings []Holding) error {
	copied := make([]Holding, len(holdings))
	now := time.Now().UTC()
	for i, holding := range holdings {
		holding.UserID = userID
		holding = holding.Normalized()
		if err := holding.Validate(); err != nil {
			return err
		}
		if holding.CreatedAt.IsZero() {
			holding.CreatedAt = now
		}
		holding.UpdatedAt = now
		copied[i] = holding
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.holdings[userID] = copied
	return nil
}

func (r *Repository) Upsert(holding Holding) (Holding, error) {
	holding = holding.Normalized()
	if err := holding.Validate(); err != nil {
		return Holding{}, err
	}
	now := time.Now().UTC()
	if holding.CreatedAt.IsZero() {
		holding.CreatedAt = now
	}
	holding.UpdatedAt = now

	r.mu.Lock()
	defer r.mu.Unlock()
	items := append([]Holding(nil), r.holdings[holding.UserID]...)
	for i, existing := range items {
		if existing.Market == holding.Market && existing.Symbol == holding.Symbol {
			items[i] = holding
			r.holdings[holding.UserID] = items
			return holding, nil
		}
	}
	r.holdings[holding.UserID] = append(items, holding)
	return holding, nil
}

func (r *Repository) Delete(userID string, market string, symbol string) bool {
	normalized := Holding{UserID: userID, Market: market, Symbol: symbol}.Normalized()
	r.mu.Lock()
	defer r.mu.Unlock()
	items := append([]Holding(nil), r.holdings[userID]...)
	for i, existing := range items {
		if existing.Market == normalized.Market && existing.Symbol == normalized.Symbol {
			r.holdings[userID] = append(items[:i], items[i+1:]...)
			return true
		}
	}
	return false
}

func (r *Repository) ListByUser(userID string) []Holding {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]Holding(nil), r.holdings[userID]...)
}
