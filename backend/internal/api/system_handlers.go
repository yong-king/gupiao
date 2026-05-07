package api

import (
	"net"
	"net/http"
	"strings"
	"time"
)

type dependencyStatus struct {
	Name       string `json:"name"`
	Configured string `json:"configured"`
	Reachable  bool   `json:"reachable"`
	Message    string `json:"message"`
}

type systemDependenciesResponse struct {
	Database    dependencyStatus `json:"database"`
	Redis       dependencyStatus `json:"redis"`
	LLM         dependencyStatus `json:"llm"`
	StockSource dependencyStatus `json:"stock_source"`
}

func (s *Server) handleSystemDependencies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "validation_error", "Method not allowed.", requestID(r))
		return
	}
	stockSource := dependencyStatus{Name: "stock_source", Message: "not configured"}
	if len(s.cfg.StockSources) > 0 {
		src := s.cfg.StockSources[0]
		stockSource = dependencyStatus{Name: src.Name, Configured: src.Type, Reachable: true, Message: src.BaseURL}
	}
	writeJSON(w, http.StatusOK, systemDependenciesResponse{
		Database:    tcpStatus("postgres", postgresAddr(s.cfg.DatabaseURL)),
		Redis:       tcpStatus("redis", s.cfg.RedisAddr),
		LLM:         dependencyStatus{Name: s.cfg.LLM.Provider, Configured: s.cfg.LLM.Model, Reachable: s.llmKeyConfigured(), Message: "api key env: " + s.cfg.LLM.APIKeyEnv},
		StockSource: stockSource,
	})
}

func tcpStatus(name string, address string) dependencyStatus {
	status := dependencyStatus{Name: name, Configured: address}
	if strings.TrimSpace(address) == "" {
		status.Message = "not configured"
		return status
	}
	conn, err := net.DialTimeout("tcp", address, 300*time.Millisecond)
	if err != nil {
		status.Message = err.Error()
		return status
	}
	_ = conn.Close()
	status.Reachable = true
	status.Message = "reachable"
	return status
}

func postgresAddr(databaseURL string) string {
	if strings.Contains(databaseURL, "@") {
		afterAt := strings.SplitN(databaseURL, "@", 2)[1]
		hostPort := strings.SplitN(afterAt, "/", 2)[0]
		if !strings.Contains(hostPort, ":") {
			return hostPort + ":5432"
		}
		return hostPort
	}
	return "localhost:5432"
}
