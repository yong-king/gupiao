package api

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"jijin/backend/internal/alerts"
	"jijin/backend/internal/delivery"
)

type deliverySettingsRequest struct {
	UserID          string `json:"user_id"`
	Enabled         bool   `json:"enabled"`
	DefaultChannel  string `json:"default_channel"`
	WeComWebhookURL string `json:"wecom_webhook_url"`
	WeComMention    string `json:"wecom_mention"`
	WeChatTarget    string `json:"wechat_target"`
	WeChatNote      string `json:"wechat_note"`
	MinRiskLevel    string `json:"min_risk_level"`
	CooldownSeconds int    `json:"cooldown_seconds"`
	MCPEnabled      bool   `json:"mcp_enabled"`
}

type deliveryTestRequest struct {
	UserID string `json:"user_id"`
}

func (s *Server) handleDeliverySettings(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		userID := strings.TrimSpace(r.URL.Query().Get("user_id"))
		if userID == "" {
			WriteError(w, http.StatusBadRequest, "validation_error", "user_id is required", requestID(r))
			return
		}
		settings := s.findDeliverySettings(userID)
		writeJSON(w, http.StatusOK, sanitizeDeliverySettings(settings))
	case http.MethodPost:
		var req deliverySettingsRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			WriteError(w, http.StatusBadRequest, "validation_error", "Invalid JSON body.", requestID(r))
			return
		}
		current := s.findDeliverySettings(req.UserID)
		webhook := strings.TrimSpace(req.WeComWebhookURL)
		if webhook == "" {
			webhook = current.WeComWebhookURL
		}
		settings := delivery.Settings{
			UserID:          req.UserID,
			Enabled:         req.Enabled,
			DefaultChannel:  delivery.ChannelType(req.DefaultChannel),
			WeComWebhookURL: webhook,
			WeComMention:    strings.TrimSpace(req.WeComMention),
			WeChatTarget:    strings.TrimSpace(req.WeChatTarget),
			WeChatNote:      strings.TrimSpace(req.WeChatNote),
			MinRiskLevel:    strings.TrimSpace(req.MinRiskLevel),
			CooldownSeconds: req.CooldownSeconds,
			MCPEnabled:      req.MCPEnabled,
		}
		if settings.MinRiskLevel == "" {
			settings.MinRiskLevel = string(alerts.RiskHigh)
		}
		if settings.DefaultChannel == "" {
			settings.DefaultChannel = delivery.ChannelTypeWeCom
		}
		if err := s.deliveries.Save(settings); err != nil {
			WriteError(w, http.StatusBadRequest, "validation_error", err.Error(), requestID(r))
			return
		}
		writeJSON(w, http.StatusOK, sanitizeDeliverySettings(settings))
	default:
		WriteError(w, http.StatusMethodNotAllowed, "validation_error", "Method not allowed.", requestID(r))
	}
}

func (s *Server) handleDeliveryTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		WriteError(w, http.StatusMethodNotAllowed, "validation_error", "Method not allowed.", requestID(r))
		return
	}
	var req deliveryTestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", "Invalid JSON body.", requestID(r))
		return
	}
	settings := s.findDeliverySettings(req.UserID)
	if !settings.Enabled {
		WriteError(w, http.StatusBadRequest, "validation_error", "外部提醒投递尚未启用。", requestID(r))
		return
	}
	result, err := s.sendDeliveryMessage(r.Context(), settings, "测试高提醒消息", "这是系统发出的测试消息，用于验证企业微信或微信类通道配置。")
	logItem := delivery.Log{
		ID:             "delivery-test-" + shortDeliveryHash(req.UserID+time.Now().UTC().Format(time.RFC3339Nano)),
		UserID:         req.UserID,
		ChannelType:    string(settings.DefaultChannel),
		ChannelTarget:  deliveryTarget(settings),
		RequestSummary: "测试高提醒消息",
		CreatedAt:      time.Now().UTC(),
	}
	if err != nil {
		settings.LastTestStatus = "failed"
		settings.LastTestMessage = err.Error()
		settings.LastTestedAt = time.Now().UTC()
		_ = s.deliveries.Save(settings)
		logItem.Status = "failed"
		logItem.ErrorMessage = err.Error()
		s.deliveryLogs.Save(logItem)
		WriteError(w, http.StatusBadGateway, "delivery_error", err.Error(), requestID(r))
		return
	}
	settings.LastTestStatus = "sent"
	settings.LastTestMessage = result
	settings.LastTestedAt = time.Now().UTC()
	_ = s.deliveries.Save(settings)
	logItem.Status = "sent"
	logItem.ResponseSummary = result
	s.deliveryLogs.Save(logItem)
	writeJSON(w, http.StatusOK, map[string]string{"status": "sent", "message": result})
}

func (s *Server) handleDeliveryLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		WriteError(w, http.StatusMethodNotAllowed, "validation_error", "Method not allowed.", requestID(r))
		return
	}
	userID := strings.TrimSpace(r.URL.Query().Get("user_id"))
	if userID == "" {
		WriteError(w, http.StatusBadRequest, "validation_error", "user_id is required", requestID(r))
		return
	}
	writeJSON(w, http.StatusOK, s.deliveryLogs.ListByUser(userID, 30))
}

func (s *Server) findDeliverySettings(userID string) delivery.Settings {
	if item, ok := s.deliveries.FindByUserID(userID); ok {
		return item
	}
	return delivery.DefaultSettings(userID)
}

func sanitizeDeliverySettings(settings delivery.Settings) delivery.Settings {
	if settings.WeComWebhookURL != "" {
		settings.WeComWebhookURL = "已配置"
	}
	return settings
}

func (s *Server) notifyExternalHighAlert(ctx context.Context, event alerts.Event) {
	// 这里只处理高等级外发。中低等级仍然只停留在系统内提醒中心。
	if event.RiskLevel != alerts.RiskHigh && event.RiskLevel != alerts.RiskCritical {
		return
	}
	settings := s.findDeliverySettings(event.UserID)
	if !settings.Enabled {
		return
	}
	if !riskAllowed(event.RiskLevel, settings.MinRiskLevel) {
		return
	}
	title := fmt.Sprintf("高提醒 %s %s:%s", event.RiskLevel, event.Market, event.Symbol)
	body := fmt.Sprintf("信号：%s\n摘要：%s\n触发时间：%s\n仅供研究提醒，不构成买卖指令。", event.Signal, event.Summary, event.DataTime.Format(time.RFC3339))
	if settings.CooldownSeconds > 0 && s.deliveryLogs.HasRecentSuccess(event.UserID, title, time.Duration(settings.CooldownSeconds)*time.Second) {
		s.deliveryLogs.Save(delivery.Log{
			ID:             "delivery-" + shortDeliveryHash(event.ID+"cooldown"+time.Now().UTC().Format(time.RFC3339Nano)),
			UserID:         event.UserID,
			AlertEventID:   event.ID,
			ChannelType:    string(settings.DefaultChannel),
			ChannelTarget:  deliveryTarget(settings),
			RequestSummary: title,
			Status:         "skipped",
			ErrorMessage:   "hit delivery cooldown window",
			CreatedAt:      time.Now().UTC(),
		})
		return
	}
	status := "sent"
	responseSummary := ""
	errMessage := ""
	response, err := s.sendDeliveryMessage(ctx, settings, title, body)
	if err != nil {
		status = "failed"
		errMessage = err.Error()
	} else {
		responseSummary = response
	}
	s.deliveryLogs.Save(delivery.Log{
		ID:              "delivery-" + shortDeliveryHash(event.ID+string(event.RiskLevel)+time.Now().UTC().Format(time.RFC3339Nano)),
		UserID:          event.UserID,
		AlertEventID:    event.ID,
		ChannelType:     string(settings.DefaultChannel),
		ChannelTarget:   deliveryTarget(settings),
		RequestSummary:  title,
		ResponseSummary: responseSummary,
		Status:          status,
		ErrorMessage:    errMessage,
		CreatedAt:       time.Now().UTC(),
	})
}

func (s *Server) sendDeliveryMessage(ctx context.Context, settings delivery.Settings, title string, body string) (string, error) {
	switch settings.DefaultChannel {
	case delivery.ChannelTypeWeChat:
		return "", fmt.Errorf("个人微信自动投递当前仅保留配置占位，请优先使用企业微信或后续合规 MCP 通道")
	case delivery.ChannelTypeWeCom, "":
		return s.wecom.Send(ctx, settings, title, body)
	default:
		return "", fmt.Errorf("unsupported delivery channel %q", settings.DefaultChannel)
	}
}

func riskAllowed(level alerts.RiskLevel, minimum string) bool {
	order := map[string]int{
		string(alerts.RiskLow):      1,
		string(alerts.RiskMedium):   2,
		string(alerts.RiskHigh):     3,
		string(alerts.RiskCritical): 4,
	}
	return order[string(level)] >= order[minimum]
}

func deliveryTarget(settings delivery.Settings) string {
	if settings.DefaultChannel == delivery.ChannelTypeWeCom {
		if settings.WeComWebhookURL != "" {
			return "wecom webhook"
		}
		return "wecom not configured"
	}
	if settings.WeChatTarget != "" {
		return settings.WeChatTarget
	}
	return string(settings.DefaultChannel)
}

func shortDeliveryHash(value string) string {
	hash := sha1.Sum([]byte(value))
	return hex.EncodeToString(hash[:])[:20]
}
