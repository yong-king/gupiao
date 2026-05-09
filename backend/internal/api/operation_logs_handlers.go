package api

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"jijin/backend/internal/persistence"
)

func persistenceOperationLog(userID string, market string, symbol string, operationType string, component string, model string, input string, output string, metadata map[string]string) persistence.OperationLog {
	return persistence.OperationLog{
		UserID:        userID,
		Market:        strings.ToUpper(strings.TrimSpace(market)),
		Symbol:        strings.ToUpper(strings.TrimSpace(symbol)),
		OperationType: operationType,
		Component:     component,
		Model:         model,
		InputSummary:  input,
		OutputSummary: output,
		Metadata:      metadata,
		CreatedAt:     time.Now().UTC(),
	}
}

func (s *Server) handleOperationLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "validation_error", "Method not allowed.", requestID(r))
		return
	}
	userID := strings.TrimSpace(r.URL.Query().Get("user_id"))
	if userID == "" {
		WriteError(w, http.StatusBadRequest, "validation_error", "user_id is required", requestID(r))
		return
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if s.store == nil {
		writeJSON(w, http.StatusOK, []persistence.OperationLog{})
		return
	}
	logs, err := s.store.ListOperationLogs(userID, limit)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "storage_error", err.Error(), requestID(r))
		return
	}
	writeJSON(w, http.StatusOK, logs)
}

func (s *Server) saveOperationLog(log persistence.OperationLog) {
	if s.store == nil {
		return
	}
	if log.ID == "" {
		log.ID = "op-" + shortHash(log.UserID+log.Market+log.Symbol+log.OperationType+time.Now().UTC().Format(time.RFC3339Nano))
	}
	// 日志不能影响主流程：AI 或行情采集成功时，日志写入失败只忽略，不中断用户操作。
	_ = s.store.SaveOperationLog(log)
}
