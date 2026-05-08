package api

import (
	"encoding/json"
	"net/http"

	"jijin/backend/internal/accounts"
	"jijin/backend/internal/settings"
)

type accountRequest struct {
	ID          string            `json:"id"`
	UserID      string            `json:"user_id"`
	Alias       string            `json:"alias"`
	RefreshMode string            `json:"refresh_mode"`
	ReadOnly    bool              `json:"read_only"`
	Metadata    map[string]string `json:"metadata"`
}

func (s *Server) handleAccounts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if s.store != nil {
			items, err := s.store.ListAccountsByUser(r.URL.Query().Get("user_id"))
			if err != nil {
				WriteError(w, http.StatusInternalServerError, "storage_error", err.Error(), requestID(r))
				return
			}
			writeJSON(w, http.StatusOK, items)
			return
		}
		writeJSON(w, http.StatusOK, s.accounts.ListByUser(r.URL.Query().Get("user_id")))
	case http.MethodPost:
		var req accountRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteError(w, http.StatusBadRequest, "validation_error", "Invalid JSON body.", requestID(r))
			return
		}
		account, err := s.accounts.Save(accounts.Config{
			ID:          req.ID,
			UserID:      req.UserID,
			Alias:       req.Alias,
			RefreshMode: settings.RefreshMode(req.RefreshMode),
			ReadOnly:    req.ReadOnly,
			Metadata:    req.Metadata,
		})
		if err != nil {
			WriteError(w, http.StatusBadRequest, "validation_error", err.Error(), requestID(r))
			return
		}
		if s.store != nil {
			if err := s.store.SaveAccount(account); err != nil {
				WriteError(w, http.StatusInternalServerError, "storage_error", err.Error(), requestID(r))
				return
			}
		}
		writeJSON(w, http.StatusCreated, account)
	default:
		WriteError(w, http.StatusMethodNotAllowed, "validation_error", "Method not allowed.", requestID(r))
	}
}
