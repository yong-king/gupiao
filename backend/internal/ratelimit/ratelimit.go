package ratelimit

import (
	"sync"
	"time"
)

type CooldownLimiter struct {
	mu       sync.Mutex
	lastSeen map[string]time.Time
	now      func() time.Time
}

func NewCooldownLimiter() *CooldownLimiter {
	return &CooldownLimiter{
		lastSeen: make(map[string]time.Time),
		now:      func() time.Time { return time.Now().UTC() },
	}
}

func NewCooldownLimiterWithClock(now func() time.Time) *CooldownLimiter {
	limiter := NewCooldownLimiter()
	limiter.now = now
	return limiter
}

func (l *CooldownLimiter) Allow(scope string, cooldown time.Duration) (bool, time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	last, ok := l.lastSeen[scope]
	if ok {
		elapsed := now.Sub(last)
		if elapsed < cooldown {
			return false, cooldown - elapsed
		}
	}
	l.lastSeen[scope] = now
	return true, 0
}

type FailureTracker struct {
	mu       sync.Mutex
	failures map[string]int
}

func NewFailureTracker() *FailureTracker {
	return &FailureTracker{failures: make(map[string]int)}
}

func (t *FailureTracker) RecordFailure(scope string) time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.failures[scope]++
	failures := t.failures[scope]
	return time.Duration(1<<min(failures-1, 5)) * time.Minute
}

func (t *FailureTracker) Reset(scope string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.failures, scope)
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
