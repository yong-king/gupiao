package api

import (
	"encoding/json"
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
