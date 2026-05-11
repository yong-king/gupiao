package delivery

import (
	"errors"
	"strings"
	"sync"
	"time"

	"jijin/backend/internal/alerts"
)

type ChannelType string

const (
	ChannelTypeWeCom  ChannelType = "wecom_webhook"
	ChannelTypeWeChat ChannelType = "wechat_placeholder"
)

type Settings struct {
	UserID                string      `json:"user_id"`
	Enabled               bool        `json:"enabled"`
	DefaultChannel        ChannelType `json:"default_channel"`
	WeComWebhookURL       string      `json:"wecom_webhook_url,omitempty"`
	WeComMention          string      `json:"wecom_mention,omitempty"`
	WeChatTarget          string      `json:"wechat_target,omitempty"`
	WeChatNote            string      `json:"wechat_note,omitempty"`
	MinRiskLevel          string      `json:"min_risk_level"`
	CooldownSeconds       int         `json:"cooldown_seconds"`
	MCPEnabled            bool        `json:"mcp_enabled"`
	LastTestStatus        string      `json:"last_test_status,omitempty"`
	LastTestMessage       string      `json:"last_test_message,omitempty"`
	LastTestedAt          time.Time   `json:"last_tested_at,omitempty"`
}

func DefaultSettings(userID string) Settings {
	return Settings{
		UserID:          userID,
		Enabled:         false,
		DefaultChannel:  ChannelTypeWeCom,
		MinRiskLevel:    string(alerts.RiskHigh),
		CooldownSeconds: 1800,
		MCPEnabled:      false,
	}
}

func (s Settings) Validate() error {
	if strings.TrimSpace(s.UserID) == "" {
		return errors.New("user id is required")
	}
	if s.DefaultChannel != "" && s.DefaultChannel != ChannelTypeWeCom && s.DefaultChannel != ChannelTypeWeChat {
		return errors.New("unsupported delivery channel")
	}
	if s.CooldownSeconds < 0 {
		return errors.New("cooldown seconds must not be negative")
	}
	if s.MinRiskLevel == "" {
		return errors.New("min risk level is required")
	}
	return nil
}

type Log struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	AlertEventID    string    `json:"alert_event_id"`
	ChannelType     string    `json:"channel_type"`
	ChannelTarget   string    `json:"channel_target"`
	RequestSummary  string    `json:"request_summary"`
	ResponseSummary string    `json:"response_summary"`
	Status          string    `json:"status"`
	ErrorMessage    string    `json:"error_message"`
	CreatedAt       time.Time `json:"created_at"`
}

type SettingsRepository struct {
	mu    sync.RWMutex
	items map[string]Settings
}

func NewSettingsRepository() *SettingsRepository {
	return &SettingsRepository{items: make(map[string]Settings)}
}

func (r *SettingsRepository) Save(settings Settings) error {
	if err := settings.Validate(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items[settings.UserID] = settings
	return nil
}

func (r *SettingsRepository) FindByUserID(userID string) (Settings, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.items[userID]
	return item, ok
}

type LogRepository struct {
	mu   sync.RWMutex
	logs []Log
}

func NewLogRepository() *LogRepository {
	return &LogRepository{}
}

func (r *LogRepository) Save(log Log) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.logs = append([]Log{log}, r.logs...)
}

func (r *LogRepository) ListByUser(userID string, limit int) []Log {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Log, 0, limit)
	for _, item := range r.logs {
		if item.UserID != userID {
			continue
		}
		out = append(out, item)
		if limit > 0 && len(out) >= limit {
			break
		}
	}
	return out
}

func (r *LogRepository) HasRecentSuccess(userID string, requestSummary string, within time.Duration) bool {
	if within <= 0 {
		return false
	}
	cutoff := time.Now().UTC().Add(-within)
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, item := range r.logs {
		if item.UserID != userID || item.Status != "sent" || item.RequestSummary != requestSummary {
			continue
		}
		if item.CreatedAt.After(cutoff) {
			return true
		}
	}
	return false
}
