package audit

import (
	"errors"
	"sync"
	"time"
)

type Entry struct {
	ID        string
	ActorID   string
	Action    string
	Target    string
	TargetID  string
	RequestID string
	Source    string
	DataTime  *time.Time
	Metadata  map[string]string
	CreatedAt time.Time
}

func (e Entry) Validate() error {
	if e.ID == "" {
		return errors.New("audit id is required")
	}
	if e.Action == "" {
		return errors.New("audit action is required")
	}
	if e.Target == "" {
		return errors.New("audit target is required")
	}
	if e.RequestID == "" {
		return errors.New("audit request id is required")
	}
	if e.Source == "" {
		return errors.New("audit source is required")
	}
	return nil
}

type Repository interface {
	Append(entry Entry) error
	List() []Entry
}

type MemoryRepository struct {
	mu      sync.RWMutex
	entries []Entry
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{}
}

func (r *MemoryRepository) Append(entry Entry) error {
	if err := entry.Validate(); err != nil {
		return err
	}
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now().UTC()
	}
	entry.Metadata = copyMap(entry.Metadata)

	r.mu.Lock()
	defer r.mu.Unlock()
	r.entries = append(r.entries, entry)
	return nil
}

func (r *MemoryRepository) List() []Entry {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]Entry, len(r.entries))
	for i, entry := range r.entries {
		entry.Metadata = copyMap(entry.Metadata)
		out[i] = entry
	}
	return out
}

func copyMap(in map[string]string) map[string]string {
	if in == nil {
		return map[string]string{}
	}
	out := make(map[string]string, len(in))
	for key, value := range in {
		out[key] = value
	}
	return out
}
