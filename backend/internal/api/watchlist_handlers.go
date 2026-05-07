package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"jijin/backend/internal/audit"
	"jijin/backend/internal/watchlist"
)

type createWatchlistRequest struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

type addSymbolRequest struct {
	Market    string  `json:"market"`
	Symbol    string  `json:"symbol"`
	Note      string  `json:"note"`
	BuyPrice  float64 `json:"buy_price"`
	SellPrice float64 `json:"sell_price"`
}

func (s *Server) handleWatchlists(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		writeJSON(w, http.StatusOK, s.watchlists.ListByUser(r.URL.Query().Get("user_id")))
		return
	}
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "validation_error", "Method not allowed.", requestID(r))
		return
	}

	var req createWatchlistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", "Invalid JSON body.", requestID(r))
		return
	}

	watchlist := watchlist.Watchlist{
		ID:     req.ID,
		UserID: req.UserID,
		Name:   req.Name,
	}
	if err := s.watchlists.Save(watchlist); err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", err.Error(), requestID(r))
		return
	}
	s.audit(audit.Entry{
		ID:        requestID(r) + ":watchlist.create",
		ActorID:   req.UserID,
		Action:    "watchlist.create",
		Target:    "watchlist",
		TargetID:  req.ID,
		RequestID: requestID(r),
		Source:    "api",
		Metadata: map[string]string{
			"name": req.Name,
		},
	})

	writeJSON(w, http.StatusCreated, watchlist)
}

func (s *Server) handleWatchlistByID(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/watchlists/")
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		WriteError(w, http.StatusNotFound, "not_found", "Watchlist not found.", requestID(r))
		return
	}

	id := parts[0]
	if len(parts) == 1 && r.Method == http.MethodGet {
		got, ok := s.watchlists.FindByID(id)
		if !ok {
			WriteError(w, http.StatusNotFound, "not_found", "Watchlist not found.", requestID(r))
			return
		}
		writeJSON(w, http.StatusOK, got)
		return
	}

	if len(parts) == 2 && parts[1] == "symbols" && r.Method == http.MethodPost {
		var req addSymbolRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteError(w, http.StatusBadRequest, "validation_error", "Invalid JSON body.", requestID(r))
			return
		}
		err := s.watchlists.AddSymbol(id, watchlist.Symbol{
			Market:    req.Market,
			Symbol:    req.Symbol,
			Note:      req.Note,
			BuyPrice:  req.BuyPrice,
			SellPrice: req.SellPrice,
		})
		if err != nil {
			WriteError(w, http.StatusBadRequest, "validation_error", err.Error(), requestID(r))
			return
		}
		got, _ := s.watchlists.FindByID(id)
		s.audit(audit.Entry{
			ID:        requestID(r) + ":watchlist.symbol.add:" + symbolKey(req),
			ActorID:   got.UserID,
			Action:    "watchlist.symbol.add",
			Target:    "watchlist",
			TargetID:  id,
			RequestID: requestID(r),
			Source:    "api",
			Metadata: map[string]string{
				"market": req.Market,
				"symbol": req.Symbol,
			},
		})
		writeJSON(w, http.StatusOK, got)
		return
	}

	WriteError(w, http.StatusNotFound, "not_found", "Watchlist endpoint not found.", requestID(r))
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func requestID(r *http.Request) string {
	if id := r.Header.Get("X-Request-ID"); id != "" {
		return id
	}
	return "unknown"
}

func (s *Server) audit(entry audit.Entry) {
	_ = s.audits.Append(entry)
}

func symbolKey(req addSymbolRequest) string {
	return strings.ToUpper(strings.TrimSpace(req.Market)) + ":" + strings.ToUpper(strings.TrimSpace(req.Symbol))
}
