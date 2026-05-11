package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Sender interface {
	Send(ctx context.Context, settings Settings, title string, body string) (string, error)
}

type WeComSender struct {
	HTTP *http.Client
}

func NewWeComSender() *WeComSender {
	return &WeComSender{HTTP: &http.Client{Timeout: 5 * time.Second}}
}

func (s *WeComSender) Send(ctx context.Context, settings Settings, title string, body string) (string, error) {
	if strings.TrimSpace(settings.WeComWebhookURL) == "" {
		return "", errors.New("wecom webhook is not configured")
	}
	content := title + "\n" + body
	if mention := strings.TrimSpace(settings.WeComMention); mention != "" {
		content += "\n@" + mention
	}
	payload := map[string]any{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"content": content,
		},
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, settings.WeComWebhookURL, bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.HTTP.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return "", fmt.Errorf("wecom webhook returned status %d", resp.StatusCode)
	}
	return "wecom webhook sent", nil
}
