package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"jijin/backend/internal/marketdata"
)

type collectMarketRequest struct {
	Market string `json:"market"`
	Symbol string `json:"symbol"`
}

type collectMarketResponse struct {
	Snapshot     marketdata.Snapshot       `json:"snapshot"`
	DailyChanges []marketdata.DailyChange  `json:"daily_changes"`
	Profile      marketdata.CompanyProfile `json:"profile"`
	Warning      string                    `json:"warning,omitempty"`
}

func (s *Server) handleMarketCollect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "validation_error", "Method not allowed.", requestID(r))
		return
	}
	var req collectMarketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", "Invalid JSON body.", requestID(r))
		return
	}
	market := strings.ToUpper(strings.TrimSpace(req.Market))
	symbol := strings.ToUpper(strings.TrimSpace(req.Symbol))
	if market == "" || symbol == "" {
		WriteError(w, http.StatusBadRequest, "validation_error", "market and symbol are required.", requestID(r))
		return
	}
	snapshots, err := s.refreshes.Provider.FetchQuotes(r.Context(), []marketdata.QuoteRequest{{Market: market, Symbol: symbol}})
	if err != nil {
		allSnapshots := s.listSnapshots(market, symbol)
		s.saveOperationLog(persistenceOperationLog("", market, symbol, "crawler_quote_collect", "market_provider", "", fmt.Sprintf("抓取实时行情 %s:%s", market, symbol), "行情源返回错误："+err.Error(), map[string]string{
			"status": "failed",
		}))
		if len(allSnapshots) > 0 {
			latest := allSnapshots[len(allSnapshots)-1]
			writeJSON(w, http.StatusOK, collectMarketResponse{
				Snapshot:     latest,
				DailyChanges: dailyChangesFromSnapshots(allSnapshots),
				Profile:      marketdata.ProfileFromSnapshots(market, symbol, allSnapshots),
				Warning:      "行情源暂时不可用，已返回最近一次保存行情。",
			})
			return
		}
		WriteError(w, http.StatusBadGateway, "provider_error", err.Error(), requestID(r))
		return
	}
	if err := s.refreshes.Snapshots.SaveAll(snapshots); err != nil {
		WriteError(w, http.StatusInternalServerError, "storage_error", err.Error(), requestID(r))
		return
	}
	if s.store != nil {
		if err := s.store.SaveSnapshots(snapshots); err != nil {
			WriteError(w, http.StatusInternalServerError, "storage_error", err.Error(), requestID(r))
			return
		}
	}
	allSnapshots := s.listSnapshots(market, symbol)
	s.saveOperationLog(persistenceOperationLog("", market, symbol, "crawler_quote_collect", "market_provider", "", fmt.Sprintf("抓取实时行情 %s:%s", market, symbol), quoteOutputSummary(snapshots[0]), map[string]string{
		"source":           snapshots[0].Source,
		"provider_request": market + ":" + symbol,
	}))
	writeJSON(w, http.StatusOK, collectMarketResponse{
		Snapshot:     snapshots[0],
		DailyChanges: dailyChangesFromSnapshots(allSnapshots),
		Profile:      marketdata.ProfileFromSnapshots(market, symbol, allSnapshots),
	})
}

func (s *Server) handleMarketSnapshots(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "validation_error", "Method not allowed.", requestID(r))
		return
	}
	market, symbol, ok := marketSymbolQuery(w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, s.listSnapshots(market, symbol))
}

func (s *Server) handleDailyChanges(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "validation_error", "Method not allowed.", requestID(r))
		return
	}
	market, symbol, ok := marketSymbolQuery(w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, dailyChangesFromSnapshots(s.listSnapshots(market, symbol)))
}

func (s *Server) handleStockProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "validation_error", "Method not allowed.", requestID(r))
		return
	}
	market, symbol, ok := marketSymbolQuery(w, r)
	if !ok {
		return
	}
	snapshots := s.listSnapshots(market, symbol)
	writeJSON(w, http.StatusOK, marketdata.ProfileFromSnapshots(market, symbol, snapshots))
}

func (s *Server) listSnapshots(market string, symbol string) []marketdata.Snapshot {
	if s.store != nil {
		if snapshots, err := s.store.ListSnapshots(market, symbol); err == nil {
			return snapshots
		}
	}
	return s.refreshes.Snapshots.ListBySymbol(market, symbol)
}

func dailyChangesFromSnapshots(snapshots []marketdata.Snapshot) []marketdata.DailyChange {
	repo := marketdata.NewSnapshotRepository()
	_ = repo.SaveAll(snapshots)
	if len(snapshots) == 0 {
		return nil
	}
	return repo.DailyChanges(snapshots[0].Market, snapshots[0].Symbol)
}

func marketSymbolQuery(w http.ResponseWriter, r *http.Request) (string, string, bool) {
	market := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("market")))
	symbol := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("symbol")))
	if market == "" || symbol == "" {
		WriteError(w, http.StatusBadRequest, "validation_error", "market and symbol are required.", requestID(r))
		return "", "", false
	}
	return market, symbol, true
}

func quoteOutputSummary(snapshot marketdata.Snapshot) string {
	return fmt.Sprintf("返回 %s:%s 最新价 %.2f，涨跌幅 %.2f%%，成交量 %d", snapshot.Market, snapshot.Symbol, snapshot.Price, snapshot.ChangePercent, snapshot.Volume)
}
