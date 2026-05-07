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
		Holdings:  s.holdings.ListByUser(req.UserID),
	})
}

func (s *Server) handleHoldings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "validation_error", "Method not allowed.", requestID(r))
		return
	}
	writeJSON(w, http.StatusOK, s.holdings.ListByUser(r.URL.Query().Get("user_id")))
}
