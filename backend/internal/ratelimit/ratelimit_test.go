package ratelimit

import (
	"testing"
	"time"
)

func TestCooldownLimiter(t *testing.T) {
	now := time.Date(2026, 5, 5, 10, 0, 0, 0, time.UTC)
	limiter := NewCooldownLimiterWithClock(func() time.Time { return now })

	allowed, wait := limiter.Allow("source:mock", time.Minute)
	if !allowed || wait != 0 {
		t.Fatalf("first request should be allowed")
	}

	allowed, wait = limiter.Allow("source:mock", time.Minute)
	if allowed || wait != time.Minute {
		t.Fatalf("second request should be rate limited for one minute, allowed=%v wait=%v", allowed, wait)
	}

	now = now.Add(time.Minute)
	allowed, wait = limiter.Allow("source:mock", time.Minute)
	if !allowed || wait != 0 {
		t.Fatalf("request after cooldown should be allowed, allowed=%v wait=%v", allowed, wait)
	}
}

func TestFailureTrackerBackoff(t *testing.T) {
	tracker := NewFailureTracker()
	if got := tracker.RecordFailure("mock"); got != time.Minute {
		t.Fatalf("unexpected first backoff %v", got)
	}
	if got := tracker.RecordFailure("mock"); got != 2*time.Minute {
		t.Fatalf("unexpected second backoff %v", got)
	}
	tracker.Reset("mock")
	if got := tracker.RecordFailure("mock"); got != time.Minute {
		t.Fatalf("unexpected reset backoff %v", got)
	}
}
