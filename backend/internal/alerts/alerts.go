package alerts

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"jijin/backend/internal/marketdata"
)

type Signal string

const (
	SignalBuyWatch        Signal = "buy_watch"
	SignalSellWatch       Signal = "sell_watch"
	SignalHoldWatch       Signal = "hold_watch"
	SignalTakeProfitWatch Signal = "take_profit_watch"
	SignalStopLossWatch   Signal = "stop_loss_watch"
	SignalRiskWarning     Signal = "risk_warning"
	SignalAbnormalMove    Signal = "abnormal_movement"
	SignalDataIssue       Signal = "data_issue"
)

type RiskLevel string

const (
	RiskLow      RiskLevel = "low"
	RiskMedium   RiskLevel = "medium"
	RiskHigh     RiskLevel = "high"
	RiskCritical RiskLevel = "critical"
)

type RuleType string

const (
	RulePriceBelow         RuleType = "price_below"
	RulePriceAbove         RuleType = "price_above"
	RuleChangePercentBelow RuleType = "change_percent_below"
	RuleChangePercentAbove RuleType = "change_percent_above"
)

type Rule struct {
	ID        string
	UserID    string
	Market    string
	Symbol    string
	Type      RuleType
	Threshold float64
	Signal    Signal
	RiskLevel RiskLevel
	Enabled   bool
	Cooldown  time.Duration
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r Rule) Validate() error {
	if r.ID == "" || r.UserID == "" || r.Market == "" || r.Symbol == "" {
		return errors.New("rule id, user id, market and symbol are required")
	}
	switch r.Type {
	case RulePriceBelow, RulePriceAbove, RuleChangePercentBelow, RuleChangePercentAbove:
	default:
		return fmt.Errorf("unsupported rule type %q", r.Type)
	}
	if r.Signal == "" {
		return errors.New("signal is required")
	}
	if r.RiskLevel == "" {
		return errors.New("risk level is required")
	}
	if r.Cooldown < 0 {
		return errors.New("cooldown must not be negative")
	}
	return nil
}

type Event struct {
	ID             string
	UserID         string
	RuleID         string
	Market         string
	Symbol         string
	Signal         Signal
	RiskLevel      RiskLevel
	TriggeredRules []string
	Summary        string
	DataTime       time.Time
	Source         string
	CreatedAt      time.Time
	Read           bool
}

type RuleRepository struct {
	mu    sync.RWMutex
	rules map[string]Rule
}

func NewRuleRepository() *RuleRepository {
	return &RuleRepository{rules: make(map[string]Rule)}
}

func (r *RuleRepository) Save(rule Rule) error {
	if err := rule.Validate(); err != nil {
		return err
	}
	now := time.Now().UTC()
	if rule.CreatedAt.IsZero() {
		rule.CreatedAt = now
	}
	rule.UpdatedAt = now
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rules[rule.ID] = rule
	return nil
}

func (r *RuleRepository) ListByUser(userID string) []Rule {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []Rule
	for _, rule := range r.rules {
		if rule.UserID == userID {
			out = append(out, rule)
		}
	}
	return out
}

func (r *RuleRepository) Delete(userID string, id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	rule, ok := r.rules[id]
	if !ok || rule.UserID != userID {
		return false
	}
	delete(r.rules, id)
	return true
}

type EventRepository struct {
	mu     sync.RWMutex
	events []Event
	now    func() time.Time
}

func NewEventRepository() *EventRepository {
	return &EventRepository{now: func() time.Time { return time.Now().UTC() }}
}

func NewEventRepositoryWithClock(now func() time.Time) *EventRepository {
	return &EventRepository{now: now}
}

func (r *EventRepository) AddIfNotDuplicate(event Event, cooldown time.Duration) (Event, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	now := r.now()
	for _, existing := range r.events {
		if existing.UserID == event.UserID && existing.RuleID == event.RuleID && existing.Signal == event.Signal && now.Sub(existing.CreatedAt) < cooldown {
			return existing, false
		}
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = now
	}
	r.events = append(r.events, event)
	return event, true
}

func (r *EventRepository) ListByUser(userID string) []Event {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []Event
	for _, event := range r.events {
		if event.UserID == userID {
			out = append(out, event)
		}
	}
	return out
}

func Evaluate(rule Rule, snapshot marketdata.Snapshot) (Event, bool) {
	if !rule.Enabled {
		return Event{}, false
	}
	if rule.Market != snapshot.Market || rule.Symbol != snapshot.Symbol {
		return Event{}, false
	}

	matched := false
	switch rule.Type {
	case RulePriceBelow:
		matched = snapshot.Price < rule.Threshold
	case RulePriceAbove:
		matched = snapshot.Price > rule.Threshold
	case RuleChangePercentBelow:
		matched = snapshot.ChangePercent < rule.Threshold
	case RuleChangePercentAbove:
		matched = snapshot.ChangePercent > rule.Threshold
	}
	if !matched {
		return Event{}, false
	}

	return Event{
		ID:             rule.ID + ":" + snapshot.DataTime.Format(time.RFC3339Nano),
		UserID:         rule.UserID,
		RuleID:         rule.ID,
		Market:         rule.Market,
		Symbol:         rule.Symbol,
		Signal:         rule.Signal,
		RiskLevel:      rule.RiskLevel,
		TriggeredRules: []string{string(rule.Type)},
		Summary:        fmt.Sprintf("%s %s triggered %s at %.2f", rule.Market, rule.Symbol, rule.Type, snapshot.Price),
		DataTime:       snapshot.DataTime,
		Source:         snapshot.Source,
	}, true
}
