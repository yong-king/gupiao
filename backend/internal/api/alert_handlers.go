package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"jijin/backend/internal/alerts"
)

type createAlertRuleRequest struct {
	ID              string  `json:"id"`
	UserID          string  `json:"user_id"`
	Market          string  `json:"market"`
	Symbol          string  `json:"symbol"`
	Type            string  `json:"type"`
	Threshold       float64 `json:"threshold"`
	Signal          string  `json:"signal"`
	RiskLevel       string  `json:"risk_level"`
	Enabled         bool    `json:"enabled"`
	CooldownSeconds int     `json:"cooldown_seconds"`
}

func (s *Server) handleAlertRules(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var req createAlertRuleRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteError(w, http.StatusBadRequest, "validation_error", "Invalid JSON body.", requestID(r))
			return
		}
		cooldown := time.Duration(req.CooldownSeconds) * time.Second
		if cooldown == 0 {
			cooldown = 30 * time.Minute
		}
		rule := alerts.Rule{
			ID:        req.ID,
			UserID:    req.UserID,
			Market:    req.Market,
			Symbol:    req.Symbol,
			Type:      alerts.RuleType(req.Type),
			Threshold: req.Threshold,
			Signal:    alerts.Signal(req.Signal),
			RiskLevel: alerts.RiskLevel(req.RiskLevel),
			Enabled:   req.Enabled,
			Cooldown:  cooldown,
		}
		if err := s.alertRules.Save(rule); err != nil {
			WriteError(w, http.StatusBadRequest, "validation_error", err.Error(), requestID(r))
			return
		}
		if s.store != nil {
			if err := s.store.SaveAlertRule(rule); err != nil {
				WriteError(w, http.StatusInternalServerError, "storage_error", err.Error(), requestID(r))
				return
			}
		}
		writeJSON(w, http.StatusCreated, rule)
	case http.MethodGet:
		writeJSON(w, http.StatusOK, s.listAlertRules(r.URL.Query().Get("user_id")))
	case http.MethodDelete:
		id := strings.TrimSpace(r.URL.Query().Get("id"))
		userID := strings.TrimSpace(r.URL.Query().Get("user_id"))
		if id == "" || userID == "" {
			WriteError(w, http.StatusBadRequest, "validation_error", "id and user_id are required", requestID(r))
			return
		}
		deleted := s.alertRules.Delete(userID, id)
		if s.store != nil {
			var err error
			deleted, err = s.store.DeleteAlertRule(userID, id)
			if err != nil {
				WriteError(w, http.StatusInternalServerError, "storage_error", err.Error(), requestID(r))
				return
			}
		}
		if !deleted {
			WriteError(w, http.StatusNotFound, "not_found", "Alert rule not found.", requestID(r))
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"deleted": true})
	default:
		WriteError(w, http.StatusMethodNotAllowed, "validation_error", "Method not allowed.", requestID(r))
	}
}

func (s *Server) handleAlerts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "validation_error", "Method not allowed.", requestID(r))
		return
	}
	if s.store != nil {
		events, err := s.store.ListAlertEventsByUser(r.URL.Query().Get("user_id"))
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "storage_error", err.Error(), requestID(r))
			return
		}
		writeJSON(w, http.StatusOK, events)
		return
	}
	writeJSON(w, http.StatusOK, s.alerts.ListByUser(r.URL.Query().Get("user_id")))
}

func (s *Server) handleNotifications(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "validation_error", "Method not allowed.", requestID(r))
		return
	}
	if s.store != nil {
		messages, err := s.store.ListNotificationsByUser(r.URL.Query().Get("user_id"))
		if err != nil {
			WriteError(w, http.StatusInternalServerError, "storage_error", err.Error(), requestID(r))
			return
		}
		writeJSON(w, http.StatusOK, messages)
		return
	}
	writeJSON(w, http.StatusOK, s.notifier.ListByUser(r.URL.Query().Get("user_id")))
}

type markNotificationReadRequest struct {
	ID string `json:"id"`
}

func (s *Server) handleNotificationRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "validation_error", "Method not allowed.", requestID(r))
		return
	}
	var req markNotificationReadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", "Invalid JSON body.", requestID(r))
		return
	}
	if s.store != nil {
		if err := s.store.MarkNotificationRead(req.ID); err != nil {
			WriteError(w, http.StatusNotFound, "not_found", err.Error(), requestID(r))
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "read"})
		return
	}
	if err := s.notifier.MarkRead(req.ID); err != nil {
		WriteError(w, http.StatusNotFound, "not_found", err.Error(), requestID(r))
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "read"})
}

func (s *Server) listAlertRules(userID string) []alerts.Rule {
	if s.store != nil {
		if rules, err := s.store.ListAlertRulesByUser(userID); err == nil {
			return rules
		}
	}
	return s.alertRules.ListByUser(userID)
}
