package alerts

import (
	"testing"
	"time"

	"jijin/backend/internal/marketdata"
)

func TestEvaluatePriceRule(t *testing.T) {
	rule := Rule{
		ID:        "rule-1",
		UserID:    "user-1",
		Market:    "US",
		Symbol:    "AAPL",
		Type:      RulePriceBelow,
		Threshold: 160,
		Signal:    SignalBuyWatch,
		RiskLevel: RiskMedium,
		Enabled:   true,
		Cooldown:  time.Minute,
	}
	snapshot := marketdata.Snapshot{Market: "US", Symbol: "AAPL", Price: 150, Source: "mock", DataTime: time.Now().UTC()}

	event, ok := Evaluate(rule, snapshot)
	if !ok {
		t.Fatal("expected rule to match")
	}
	if event.Signal != SignalBuyWatch || event.RiskLevel != RiskMedium {
		t.Fatalf("unexpected event: %#v", event)
	}
}

func TestEvaluateDisabledRuleDoesNotTrigger(t *testing.T) {
	rule := Rule{Enabled: false, Type: RulePriceBelow}
	_, ok := Evaluate(rule, marketdata.Snapshot{})
	if ok {
		t.Fatal("disabled rule should not trigger")
	}
}

func TestEventRepositoryDeduplicatesWithinCooldown(t *testing.T) {
	now := time.Date(2026, 5, 5, 10, 0, 0, 0, time.UTC)
	repo := NewEventRepositoryWithClock(func() time.Time { return now })
	event := Event{ID: "event-1", UserID: "user-1", RuleID: "rule-1", Signal: SignalRiskWarning}

	if _, created := repo.AddIfNotDuplicate(event, time.Minute); !created {
		t.Fatal("first event should be created")
	}
	if _, created := repo.AddIfNotDuplicate(event, time.Minute); created {
		t.Fatal("duplicate event should be suppressed")
	}

	now = now.Add(time.Minute)
	if _, created := repo.AddIfNotDuplicate(Event{ID: "event-2", UserID: "user-1", RuleID: "rule-1", Signal: SignalRiskWarning}, time.Minute); !created {
		t.Fatal("event after cooldown should be created")
	}
}
