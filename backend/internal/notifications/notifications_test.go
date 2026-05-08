package notifications

import (
	"testing"
	"time"

	"jijin/backend/internal/alerts"
)

func TestPublishListAndMarkRead(t *testing.T) {
	center := NewCenter()
	message := FromEvent(alerts.Event{
		ID:        "event-1",
		UserID:    "user-1",
		Market:    "US",
		Symbol:    "AAPL",
		Signal:    alerts.SignalRiskWarning,
		RiskLevel: alerts.RiskHigh,
		Summary:   "Risk increased.",
		DataTime:  time.Now().UTC(),
	})

	if err := center.Publish(message); err != nil {
		t.Fatalf("publish message: %v", err)
	}
	got := center.ListByUser("user-1")
	if len(got) != 1 || got[0].Title == "" {
		t.Fatalf("unexpected messages: %#v", got)
	}
	if err := center.MarkRead("event-1"); err != nil {
		t.Fatalf("mark read: %v", err)
	}
	got = center.ListByUser("user-1")
	if !got[0].Read {
		t.Fatal("expected message to be read")
	}
}

func TestMarkReadRejectsMissingMessage(t *testing.T) {
	center := NewCenter()
	if err := center.MarkRead("missing"); err == nil {
		t.Fatal("expected missing message to fail")
	}
}
