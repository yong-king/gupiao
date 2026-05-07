package api

import (
	"encoding/json"
	"net/http"
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
		writeJSON(w, http.StatusCreated, rule)
	case http.MethodGet:
		writeJSON(w, http.StatusOK, s.alertRules.ListByUser(r.URL.Query().Get("user_id")))
	default:
		WriteError(w, http.StatusMethodNotAllowed, "validation_error", "Method not allowed.", requestID(r))
	}
}

func (s *Server) handleAlerts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "validation_error", "Method not allowed.", requestID(r))
		return
	}
	writeJSON(w, http.StatusOK, s.alerts.ListByUser(r.URL.Query().Get("user_id")))
}

func (s *Server) handleNotifications(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "validation_error", "Method not allowed.", requestID(r))
		return
	}
	writeJSON(w, http.StatusOK, s.notifier.ListByUser(r.URL.Query().Get("user_id")))
}
