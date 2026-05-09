package marketdata

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type StooqProvider struct {
	BaseURL string
	Client  *http.Client
}

func NewStooqProvider(baseURL string) *StooqProvider {
	if strings.TrimSpace(baseURL) == "" {
		baseURL = "https://stooq.com/q/l/"
	}
	return &StooqProvider{
		BaseURL: baseURL,
		Client:  &http.Client{Timeout: 8 * time.Second},
	}
}

func (p *StooqProvider) FetchQuotes(ctx context.Context, requests []QuoteRequest) ([]Snapshot, error) {
	snapshots := make([]Snapshot, 0, len(requests))
	for _, request := range requests {
		snapshot, err := p.FetchQuote(ctx, request)
		if err != nil {
			return nil, err
		}
		snapshots = append(snapshots, snapshot)
	}
	return snapshots, nil
}

func (p *StooqProvider) FetchQuote(ctx context.Context, request QuoteRequest) (Snapshot, error) {
	if strings.EqualFold(strings.TrimSpace(request.Market), "CN") {
		snapshot, err := FetchEastmoneyQuote(ctx, request, p.client())
		if err == nil {
			return snapshot, nil
		}
		fallback, fallbackErr := FetchTencentQuote(ctx, request, p.client())
		if fallbackErr != nil {
			return Snapshot{}, fmt.Errorf("eastmoney failed: %v; tencent failed: %w", err, fallbackErr)
		}
		return fallback, nil
	}
	u, err := url.Parse(p.BaseURL)
	if err != nil {
		return Snapshot{}, err
	}
	query := u.Query()
	query.Set("s", stooqSymbol(request.Market, request.Symbol))
	query.Set("f", "snd2t2ohlcv")
	query.Set("h", "")
	query.Set("e", "csv")
	u.RawQuery = query.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return Snapshot{}, err
	}
	resp, err := p.client().Do(httpReq)
	if err != nil {
		return Snapshot{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Snapshot{}, fmt.Errorf("stooq returned status %d", resp.StatusCode)
	}
	return ParseStooqCSV(resp.Body, request)
}

func ParseStooqCSV(reader io.Reader, request QuoteRequest) (Snapshot, error) {
	records, err := csv.NewReader(reader).ReadAll()
	if err != nil {
		return Snapshot{}, err
	}
	if len(records) < 2 {
		return Snapshot{}, errors.New("stooq response missing quote row")
	}
	header := records[0]
	row := records[1]
	value := func(name string) string {
		for i, column := range header {
			if strings.EqualFold(column, name) && i < len(row) {
				return strings.TrimSpace(row[i])
			}
		}
		return ""
	}
	closePrice, err := parseFloat(value("Close"))
	if err != nil || closePrice == 0 {
		return Snapshot{}, fmt.Errorf("stooq quote missing close price for %s", request.Key())
	}
	open, _ := parseFloat(value("Open"))
	high, _ := parseFloat(value("High"))
	low, _ := parseFloat(value("Low"))
	volume, _ := strconv.ParseInt(value("Volume"), 10, 64)
	dataTime := time.Now().UTC()
	if date := value("Date"); date != "" {
		clock := value("Time")
		if clock == "" {
			clock = "00:00:00"
		}
		if parsed, err := time.Parse("2006-01-02 15:04:05", date+" "+clock); err == nil {
			dataTime = parsed.UTC()
		}
	}
	baseline := open
	changePercent := 0.0
	if baseline != 0 {
		changePercent = (closePrice - baseline) / baseline * 100
	}
	return Snapshot{
		Market:        strings.ToUpper(strings.TrimSpace(request.Market)),
		Symbol:        strings.ToUpper(strings.TrimSpace(request.Symbol)),
		Name:          value("Name"),
		Open:          open,
		High:          high,
		Low:           low,
		Price:         closePrice,
		PreviousClose: baseline,
		ChangePercent: changePercent,
		Volume:        volume,
		Source:        "stooq",
		DataTime:      dataTime,
		CreatedAt:     time.Now().UTC(),
	}, nil
}

func stooqSymbol(market string, symbol string) string {
	normalized := strings.ToLower(strings.TrimSpace(symbol))
	switch strings.ToUpper(strings.TrimSpace(market)) {
	case "US":
		if !strings.Contains(normalized, ".") {
			return normalized + ".us"
		}
	case "HK":
		if !strings.Contains(normalized, ".") {
			return normalized + ".hk"
		}
	case "CN":
		if !strings.Contains(normalized, ".") {
			return normalized + ".cn"
		}
	}
	return normalized
}

func (p *StooqProvider) client() *http.Client {
	if p.Client != nil {
		return p.Client
	}
	return &http.Client{Timeout: 8 * time.Second}
}

func parseFloat(value string) (float64, error) {
	value = strings.TrimSpace(value)
	if value == "" || strings.EqualFold(value, "N/D") {
		return 0, errors.New("empty numeric value")
	}
	return strconv.ParseFloat(value, 64)
}
