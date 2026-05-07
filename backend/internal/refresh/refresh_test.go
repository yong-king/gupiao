package refresh

import (
	"context"
	"testing"
	"time"

	"jijin/backend/internal/marketdata"
	"jijin/backend/internal/ratelimit"
	"jijin/backend/internal/watchlist"
)

func TestRefreshWatchlistSucceeds(t *testing.T) {
	provider := marketdata.NewMockProvider()
	provider.SetQuote(marketdata.Snapshot{Market: "US", Symbol: "AAPL", Price: 180})
	service := NewService(provider, marketdata.NewSnapshotRepository(), NewJobRepository())
	service.Cooldown = time.Minute

	wl := watchlist.Watchlist{ID: "wl-1", UserID: "user-1", Name: "Core", Symbols: []watchlist.Symbol{{Market: "US", Symbol: "AAPL"}}}
	job, err := service.RefreshWatchlist(context.Background(), "job-1", "user-1", wl)
	if err != nil {
		t.Fatalf("refresh watchlist: %v", err)
	}
	if job.Status != StatusSucceeded || len(job.Snapshots) != 1 {
		t.Fatalf("unexpected job: %#v", job)
	}
}

func TestRefreshWatchlistFailure(t *testing.T) {
	provider := marketdata.NewMockProvider()
	service := NewService(provider, marketdata.NewSnapshotRepository(), NewJobRepository())
	wl := watchlist.Watchlist{ID: "wl-1", UserID: "user-1", Name: "Core", Symbols: []watchlist.Symbol{{Market: "US", Symbol: "AAPL"}}}

	job, err := service.RefreshWatchlist(context.Background(), "job-1", "user-1", wl)
	if err == nil {
		t.Fatal("expected refresh failure")
	}
	if job.Status != StatusFailed || job.Error == "" {
		t.Fatalf("unexpected failed job: %#v", job)
	}
}

func TestRefreshWatchlistRateLimited(t *testing.T) {
	now := time.Date(2026, 5, 5, 10, 0, 0, 0, time.UTC)
	provider := marketdata.NewMockProvider()
	provider.SetQuote(marketdata.Snapshot{Market: "US", Symbol: "AAPL", Price: 180})
	service := NewService(provider, marketdata.NewSnapshotRepository(), NewJobRepository())
	service.Cooldown = time.Minute
	service.Now = func() time.Time { return now }
	service.Limiter = ratelimit.NewCooldownLimiterWithClock(func() time.Time { return now })

	wl := watchlist.Watchlist{ID: "wl-1", UserID: "user-1", Name: "Core", Symbols: []watchlist.Symbol{{Market: "US", Symbol: "AAPL"}}}
	if _, err := service.RefreshWatchlist(context.Background(), "job-1", "user-1", wl); err != nil {
		t.Fatalf("first refresh should succeed: %v", err)
	}

	job, err := service.RefreshWatchlist(context.Background(), "job-2", "user-1", wl)
	if err != nil {
		t.Fatalf("rate limited refresh should not return transport error: %v", err)
	}
	if job.Status != StatusRateLimited {
		t.Fatalf("expected rate limited, got %#v", job)
	}
}

func TestSchedulerRunOnce(t *testing.T) {
	provider := marketdata.NewMockProvider()
	provider.SetQuote(marketdata.Snapshot{Market: "US", Symbol: "AAPL", Price: 180})
	service := NewService(provider, marketdata.NewSnapshotRepository(), NewJobRepository())
	wl := watchlist.Watchlist{ID: "wl-1", UserID: "user-1", Name: "Core", Symbols: []watchlist.Symbol{{Market: "US", Symbol: "AAPL"}}}

	job, err := (Scheduler{Service: service}).RunOnce(context.Background(), "job-1", "user-1", wl)
	if err != nil {
		t.Fatalf("scheduler run once: %v", err)
	}
	if job.Status != StatusSucceeded {
		t.Fatalf("unexpected job status %q", job.Status)
	}
}
