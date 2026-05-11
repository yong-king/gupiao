package api

import (
	"context"
	"encoding/json"
	"net/http"

	"jijin/backend/internal/alerts"
	"jijin/backend/internal/audit"
	"jijin/backend/internal/marketdata"
	"jijin/backend/internal/notifications"
)

type manualRefreshRequest struct {
	UserID      string `json:"user_id"`
	WatchlistID string `json:"watchlist_id"`
	JobID       string `json:"job_id"`
}

func (s *Server) handleManualRefresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "validation_error", "Method not allowed.", requestID(r))
		return
	}

	var req manualRefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", "Invalid JSON body.", requestID(r))
		return
	}
	if req.JobID == "" {
		req.JobID = requestID(r) + ":refresh"
	}

	wl, ok := s.watchlists.FindByID(req.WatchlistID)
	if !ok && s.store != nil {
		wl, ok = s.store.FindWatchlistByID(req.WatchlistID)
	}
	if !ok {
		WriteError(w, http.StatusNotFound, "not_found", "Watchlist not found.", requestID(r))
		return
	}
	if wl.UserID != req.UserID {
		WriteError(w, http.StatusBadRequest, "validation_error", "watchlist does not belong to user", requestID(r))
		return
	}

	job, err := s.refreshes.RefreshWatchlist(r.Context(), req.JobID, req.UserID, wl)
	if err != nil {
		s.audit(audit.Entry{
			ID:        requestID(r) + ":refresh.failed",
			ActorID:   req.UserID,
			Action:    "refresh.failed",
			Target:    "watchlist",
			TargetID:  req.WatchlistID,
			RequestID: requestID(r),
			Source:    "api",
			Metadata: map[string]string{
				"error": err.Error(),
			},
		})
		writeJSON(w, http.StatusOK, job)
		return
	}

	s.audit(audit.Entry{
		ID:        requestID(r) + ":refresh.manual",
		ActorID:   req.UserID,
		Action:    "refresh.manual",
		Target:    "watchlist",
		TargetID:  req.WatchlistID,
		RequestID: requestID(r),
		Source:    "api",
		Metadata: map[string]string{
			"status": string(job.Status),
		},
	})
	if s.store != nil {
		_ = s.store.SaveSnapshots(job.Snapshots)
	}
	s.evaluateAlerts(r.Context(), req.UserID, job.Snapshots)
	writeJSON(w, http.StatusOK, job)
}

func (s *Server) evaluateAlerts(ctx context.Context, userID string, snapshots []marketdata.Snapshot) {
	rules := s.listAlertRules(userID)
	for _, snapshot := range snapshots {
		for _, rule := range rules {
			event, ok := alerts.Evaluate(rule, snapshot)
			if !ok {
				continue
			}
			createdEvent, created := s.alerts.AddIfNotDuplicate(event, rule.Cooldown)
			if !created {
				continue
			}
			message := notifications.FromEvent(createdEvent)
			if s.store != nil {
				_ = s.store.SaveAlertEvent(createdEvent)
				_ = s.store.SaveNotification(message)
			}
			_ = s.notifier.Publish(message)
			s.notifyExternalHighAlert(ctx, createdEvent)
		}
	}
}
