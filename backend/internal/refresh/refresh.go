package refresh

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"jijin/backend/internal/marketdata"
	"jijin/backend/internal/ratelimit"
	"jijin/backend/internal/watchlist"
)

type Status string

const (
	StatusPending     Status = "pending"
	StatusRunning     Status = "running"
	StatusSucceeded   Status = "succeeded"
	StatusFailed      Status = "failed"
	StatusRateLimited Status = "rate_limited"
)

type Job struct {
	ID          string
	UserID      string
	Scope       string
	Status      Status
	Error       string
	Snapshots   []marketdata.Snapshot
	RequestedAt time.Time
	StartedAt   time.Time
	FinishedAt  time.Time
}

type JobRepository struct {
	mu   sync.RWMutex
	jobs map[string]Job
}

func NewJobRepository() *JobRepository {
	return &JobRepository{jobs: make(map[string]Job)}
}

func (r *JobRepository) Save(job Job) error {
	if job.ID == "" {
		return errors.New("job id is required")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.jobs[job.ID] = copyJob(job)
	return nil
}

func (r *JobRepository) FindByID(id string) (Job, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	job, ok := r.jobs[id]
	return copyJob(job), ok
}

func copyJob(job Job) Job {
	job.Snapshots = append([]marketdata.Snapshot(nil), job.Snapshots...)
	return job
}

type Service struct {
	Provider  marketdata.Provider
	Snapshots *marketdata.SnapshotRepository
	Jobs      *JobRepository
	Limiter   *ratelimit.CooldownLimiter
	Failures  *ratelimit.FailureTracker
	Cooldown  time.Duration
	Now       func() time.Time
}

func NewService(provider marketdata.Provider, snapshots *marketdata.SnapshotRepository, jobs *JobRepository) *Service {
	return &Service{
		Provider:  provider,
		Snapshots: snapshots,
		Jobs:      jobs,
		Limiter:   ratelimit.NewCooldownLimiter(),
		Failures:  ratelimit.NewFailureTracker(),
		Cooldown:  15 * time.Minute,
		Now:       func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) RefreshWatchlist(ctx context.Context, jobID string, userID string, wl watchlist.Watchlist) (Job, error) {
	scope := "watchlist:" + wl.ID
	now := s.Now()
	job := Job{
		ID:          jobID,
		UserID:      userID,
		Scope:       scope,
		Status:      StatusPending,
		RequestedAt: now,
	}

	if allowed, wait := s.Limiter.Allow(scope, s.Cooldown); !allowed {
		job.Status = StatusRateLimited
		job.Error = fmt.Sprintf("refresh cooldown active: retry after %s", wait)
		job.FinishedAt = now
		_ = s.Jobs.Save(job)
		return job, nil
	}

	job.Status = StatusRunning
	job.StartedAt = now
	_ = s.Jobs.Save(job)

	requests := make([]marketdata.QuoteRequest, len(wl.Symbols))
	for i, symbol := range wl.Symbols {
		requests[i] = marketdata.FromWatchlistSymbol(symbol)
	}
	if len(requests) == 0 {
		job.Status = StatusFailed
		job.Error = "watchlist has no symbols"
		job.FinishedAt = s.Now()
		_ = s.Jobs.Save(job)
		return job, errors.New(job.Error)
	}

	snapshots, err := s.Provider.FetchQuotes(ctx, requests)
	if err != nil {
		job.Status = StatusFailed
		job.Error = err.Error()
		job.FinishedAt = s.Now()
		s.Failures.RecordFailure(scope)
		_ = s.Jobs.Save(job)
		return job, err
	}

	if err := s.Snapshots.SaveAll(snapshots); err != nil {
		job.Status = StatusFailed
		job.Error = err.Error()
		job.FinishedAt = s.Now()
		_ = s.Jobs.Save(job)
		return job, err
	}

	s.Failures.Reset(scope)
	job.Status = StatusSucceeded
	job.Snapshots = snapshots
	job.FinishedAt = s.Now()
	_ = s.Jobs.Save(job)
	return job, nil
}

type Scheduler struct {
	Service *Service
}

func (s Scheduler) RunOnce(ctx context.Context, jobID string, userID string, wl watchlist.Watchlist) (Job, error) {
	return s.Service.RefreshWatchlist(ctx, jobID, userID, wl)
}
