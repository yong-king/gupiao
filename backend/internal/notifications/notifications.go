package notifications

import (
	"errors"
	"sync"
	"time"

	"jijin/backend/internal/alerts"
)

type Message struct {
	ID        string
	UserID    string
	Title     string
	Summary   string
	Signal    alerts.Signal
	RiskLevel alerts.RiskLevel
	Market    string
	Symbol    string
	DataTime  time.Time
	Read      bool
	CreatedAt time.Time
}

type Center struct {
	mu       sync.RWMutex
	messages map[string]Message
	now      func() time.Time
}

func NewCenter() *Center {
	return &Center{messages: make(map[string]Message), now: func() time.Time { return time.Now().UTC() }}
}

func FromEvent(event alerts.Event) Message {
	return Message{
		ID:        event.ID,
		UserID:    event.UserID,
		Title:     string(event.Signal) + " " + event.Market + ":" + event.Symbol,
		Summary:   event.Summary,
		Signal:    event.Signal,
		RiskLevel: event.RiskLevel,
		Market:    event.Market,
		Symbol:    event.Symbol,
		DataTime:  event.DataTime,
	}
}

func (c *Center) Publish(message Message) error {
	if message.ID == "" || message.UserID == "" {
		return errors.New("message id and user id are required")
	}
	if message.CreatedAt.IsZero() {
		message.CreatedAt = c.now()
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.messages[message.ID] = message
	return nil
}

func (c *Center) ListByUser(userID string) []Message {
	c.mu.RLock()
	defer c.mu.RUnlock()
	var out []Message
	for _, message := range c.messages {
		if message.UserID == userID {
			out = append(out, message)
		}
	}
	return out
}

func (c *Center) MarkRead(id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	message, ok := c.messages[id]
	if !ok {
		return errors.New("message not found")
	}
	message.Read = true
	c.messages[id] = message
	return nil
}
