package stocksource

import (
	"context"
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"time"

	"jijin/backend/internal/marketdata"
)

type Collector interface {
	Collect(ctx context.Context, market string, symbol string) (marketdata.Snapshot, error)
}

type PlatformPayload struct {
	Market        string  `json:"market"`
	Symbol        string  `json:"symbol"`
	Price         float64 `json:"price"`
	PreviousClose float64 `json:"previous_close"`
	ChangePercent float64 `json:"change_percent"`
	Volume        int64   `json:"volume"`
	Source        string  `json:"source"`
	DataTime      string  `json:"data_time"`
}

func ParseJSON(data []byte) (marketdata.Snapshot, error) {
	var payload PlatformPayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return marketdata.Snapshot{}, err
	}
	return payload.snapshot()
}

func ParseHTML(data string) (marketdata.Snapshot, error) {
	value := func(name string) string {
		re := regexp.MustCompile(`data-` + name + `=["']([^"']+)["']`)
		match := re.FindStringSubmatch(data)
		if len(match) != 2 {
			return ""
		}
		return match[1]
	}
	price, err := strconv.ParseFloat(value("price"), 64)
	if err != nil {
		return marketdata.Snapshot{}, errors.New("missing or invalid price")
	}
	prev, _ := strconv.ParseFloat(value("previous-close"), 64)
	change, _ := strconv.ParseFloat(value("change-percent"), 64)
	volume, _ := strconv.ParseInt(value("volume"), 10, 64)
	payload := PlatformPayload{
		Market:        value("market"),
		Symbol:        value("symbol"),
		Price:         price,
		PreviousClose: prev,
		ChangePercent: change,
		Volume:        volume,
		Source:        value("source"),
		DataTime:      value("data-time"),
	}
	return payload.snapshot()
}

func (p PlatformPayload) snapshot() (marketdata.Snapshot, error) {
	if p.Market == "" || p.Symbol == "" || p.Price == 0 {
		return marketdata.Snapshot{}, errors.New("market, symbol and price are required")
	}
	dataTime := time.Now().UTC()
	if p.DataTime != "" {
		parsed, err := time.Parse(time.RFC3339, p.DataTime)
		if err != nil {
			return marketdata.Snapshot{}, err
		}
		dataTime = parsed
	}
	source := p.Source
	if source == "" {
		source = "stock_platform"
	}
	return marketdata.Snapshot{
		Market:        p.Market,
		Symbol:        p.Symbol,
		Price:         p.Price,
		PreviousClose: p.PreviousClose,
		ChangePercent: p.ChangePercent,
		Volume:        p.Volume,
		Source:        source,
		DataTime:      dataTime,
	}, nil
}

type MockCollector struct {
	Snapshot marketdata.Snapshot
	Err      error
}

func (c MockCollector) Collect(ctx context.Context, market string, symbol string) (marketdata.Snapshot, error) {
	if c.Err != nil {
		return marketdata.Snapshot{}, c.Err
	}
	if c.Snapshot.Market == "" {
		return marketdata.Snapshot{}, errors.New("mock snapshot missing")
	}
	return c.Snapshot, nil
}
