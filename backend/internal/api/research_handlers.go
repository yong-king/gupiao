package api

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"jijin/backend/internal/holdings"
	"jijin/backend/internal/marketdata"
	"jijin/backend/internal/persistence"
)

type collectResearchRequest struct {
	UserID         string `json:"user_id"`
	Market         string `json:"market"`
	Symbol         string `json:"symbol"`
	AttentionLevel string `json:"attention_level"`
}

type collectResearchResponse struct {
	DocumentID      string                    `json:"document_id"`
	Market          string                    `json:"market"`
	Symbol          string                    `json:"symbol"`
	AttentionLevel  string                    `json:"attention_level"`
	RefreshInterval string                    `json:"refresh_interval"`
	Summary         string                    `json:"summary"`
	Profile         marketdata.CompanyProfile `json:"profile"`
	Snapshot        marketdata.Snapshot       `json:"snapshot"`
	Metadata        map[string]string         `json:"metadata"`
}

func (s *Server) handleResearchCollect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "validation_error", "Method not allowed.", requestID(r))
		return
	}
	var req collectResearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", "Invalid JSON body.", requestID(r))
		return
	}
	userID := strings.TrimSpace(req.UserID)
	market := strings.ToUpper(strings.TrimSpace(req.Market))
	symbol := strings.ToUpper(strings.TrimSpace(req.Symbol))
	if userID == "" || market == "" || symbol == "" {
		WriteError(w, http.StatusBadRequest, "validation_error", "user_id, market and symbol are required", requestID(r))
		return
	}
	attention := holdings.NormalizeAttentionLevel(req.AttentionLevel)
	snapshots := s.listSnapshots(market, symbol)
	profile := marketdata.ProfileFromSnapshots(market, symbol, snapshots)
	var snapshot marketdata.Snapshot
	if len(snapshots) > 0 {
		snapshot = snapshots[len(snapshots)-1]
	}
	interval := holdings.AttentionRefreshInterval(attention)
	summary := researchSummary(profile, snapshot, attention, interval)
	metadata := map[string]string{
		"attention_level":  attention,
		"refresh_interval": interval.String(),
		"embedding_status": "pending_embedding",
		"source":           "local_profile_price_summary",
	}
	id := ragDocumentID(userID, market, symbol, time.Now().UTC())
	if s.store != nil {
		if err := s.store.SaveRAGDocument(persistence.RAGDocument{
			ID:         id,
			UserID:     userID,
			Market:     market,
			Symbol:     symbol,
			SourceType: "stock_research_summary",
			SourceID:   market + ":" + symbol,
			Content:    summary,
			Metadata:   metadata,
			CreatedAt:  time.Now().UTC(),
		}); err != nil {
			WriteError(w, http.StatusInternalServerError, "storage_error", err.Error(), requestID(r))
			return
		}
	}
	writeJSON(w, http.StatusOK, collectResearchResponse{
		DocumentID:      id,
		Market:          market,
		Symbol:          symbol,
		AttentionLevel:  attention,
		RefreshInterval: interval.String(),
		Summary:         summary,
		Profile:         profile,
		Snapshot:        snapshot,
		Metadata:        metadata,
	})
}

func researchSummary(profile marketdata.CompanyProfile, snapshot marketdata.Snapshot, attention string, interval time.Duration) string {
	change := "暂无已保存涨跌数据"
	if snapshot.Symbol != "" {
		change = fmt.Sprintf("最新价 %.2f，涨跌幅 %.2f%%，数据源 %s", snapshot.Price, snapshot.ChangePercent, snapshot.Source)
	}
	return fmt.Sprintf("%s:%s 关注等级为%s，建议公开信息采集周期为%s。公司/产品：%s；业务：%s；行情：%s。综合建议：%s。仅用于研究和提醒，不构成投资建议。",
		profile.Market, profile.Symbol, attention, interval, strings.Join(profile.Products, "、"), profile.Business, change, profile.Analysis)
}

func ragDocumentID(userID string, market string, symbol string, at time.Time) string {
	hash := sha1.Sum([]byte(userID + ":" + market + ":" + symbol + ":" + at.Format(time.RFC3339Nano)))
	return "rag-" + hex.EncodeToString(hash[:])[:20]
}
