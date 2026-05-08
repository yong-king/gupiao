package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"jijin/backend/internal/audit"
	"jijin/backend/internal/holdings"
)

type importHoldingsRequest struct {
	UserID string `json:"user_id"`
	CSV    string `json:"csv"`
}

type upsertHoldingRequest struct {
	UserID    string  `json:"user_id"`
	Market    string  `json:"market"`
	Symbol    string  `json:"symbol"`
	Quantity  float64 `json:"quantity"`
	CostBasis float64 `json:"cost_basis"`
}

type importHoldingsResponse struct {
	Imported  int                `json:"imported"`
	RowErrors []string           `json:"row_errors"`
	Holdings  []holdings.Holding `json:"holdings"`
}

func (s *Server) handleHoldingsImport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "validation_error", "Method not allowed.", requestID(r))
		return
	}

	var req importHoldingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", "Invalid JSON body.", requestID(r))
		return
	}
	if strings.TrimSpace(req.UserID) == "" {
		WriteError(w, http.StatusBadRequest, "validation_error", "user_id is required", requestID(r))
		return
	}

	parsed, rowErrors, err := holdings.ParseCSV(req.UserID, strings.NewReader(req.CSV))
	if err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", err.Error(), requestID(r))
		return
	}
	if err := s.holdings.ReplaceForUser(req.UserID, parsed); err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", err.Error(), requestID(r))
		return
	}
	if s.store != nil {
		if err := s.store.ReplaceHoldingsForUser(req.UserID, parsed); err != nil {
			WriteError(w, http.StatusInternalServerError, "storage_error", err.Error(), requestID(r))
			return
		}
	}
	s.audit(audit.Entry{
		ID:        requestID(r) + ":holdings.import",
		ActorID:   req.UserID,
		Action:    "holdings.import",
		Target:    "holdings",
		TargetID:  req.UserID,
		RequestID: requestID(r),
		Source:    "api",
		Metadata: map[string]string{
			"imported":   strconv.Itoa(len(parsed)),
			"row_errors": strconv.Itoa(len(rowErrors)),
		},
	})

	errors := make([]string, len(rowErrors))
	for i, rowError := range rowErrors {
		errors[i] = rowError.Error()
	}

	writeJSON(w, http.StatusOK, importHoldingsResponse{
		Imported:  len(parsed),
		RowErrors: errors,
		Holdings:  s.listHoldings(req.UserID),
	})
}

func (s *Server) handleHoldings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, s.listHoldings(r.URL.Query().Get("user_id")))
	case http.MethodPost:
		var req upsertHoldingRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteError(w, http.StatusBadRequest, "validation_error", "Invalid JSON body.", requestID(r))
			return
		}
		input := holdings.Holding{
			UserID:    req.UserID,
			Market:    req.Market,
			Symbol:    req.Symbol,
			Quantity:  req.Quantity,
			CostBasis: req.CostBasis,
		}
		holding, err := s.holdings.Upsert(input)
		if err != nil {
			WriteError(w, http.StatusBadRequest, "validation_error", err.Error(), requestID(r))
			return
		}
		if s.store != nil {
			holding, err = s.store.UpsertHolding(input)
			if err != nil {
				WriteError(w, http.StatusInternalServerError, "storage_error", err.Error(), requestID(r))
				return
			}
		}
		s.audit(audit.Entry{
			ID:        requestID(r) + ":holdings.upsert",
			ActorID:   req.UserID,
			Action:    "holdings.upsert",
			Target:    "holding",
			TargetID:  holding.Market + ":" + holding.Symbol,
			RequestID: requestID(r),
			Source:    "api",
		})
		writeJSON(w, http.StatusOK, holding)
	case http.MethodDelete:
		userID := strings.TrimSpace(r.URL.Query().Get("user_id"))
		market := strings.TrimSpace(r.URL.Query().Get("market"))
		symbol := strings.TrimSpace(r.URL.Query().Get("symbol"))
		if userID == "" || market == "" || symbol == "" {
			WriteError(w, http.StatusBadRequest, "validation_error", "user_id, market and symbol are required", requestID(r))
			return
		}
		deleted := s.holdings.Delete(userID, market, symbol)
		if s.store != nil {
			var err error
			deleted, err = s.store.DeleteHolding(userID, market, symbol)
			if err != nil {
				WriteError(w, http.StatusInternalServerError, "storage_error", err.Error(), requestID(r))
				return
			}
		}
		if !deleted {
			WriteError(w, http.StatusNotFound, "not_found", "Holding not found.", requestID(r))
			return
		}
		s.audit(audit.Entry{
			ID:        requestID(r) + ":holdings.delete",
			ActorID:   userID,
			Action:    "holdings.delete",
			Target:    "holding",
			TargetID:  strings.ToUpper(market) + ":" + strings.ToUpper(symbol),
			RequestID: requestID(r),
			Source:    "api",
		})
		writeJSON(w, http.StatusOK, map[string]any{"deleted": true})
	default:
		WriteError(w, http.StatusMethodNotAllowed, "validation_error", "Method not allowed.", requestID(r))
	}
}

func (s *Server) listHoldings(userID string) []holdings.Holding {
	if s.store != nil {
		if items, err := s.store.ListHoldingsByUser(userID); err == nil {
			return items
		}
	}
	return s.holdings.ListByUser(userID)
}
