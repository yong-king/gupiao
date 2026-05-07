package observability

import "sync"

type LogEvent struct {
	Service   string `json:"service"`
	Action    string `json:"action"`
	Status    string `json:"status"`
	RequestID string `json:"request_id"`
	Message   string `json:"message"`
}

func (e LogEvent) Valid() bool {
	return e.Service != "" && e.Action != "" && e.Status != "" && e.RequestID != ""
}

type Counters struct {
	mu     sync.RWMutex
	values map[string]int
}

func NewCounters() *Counters {
	return &Counters{values: make(map[string]int)}
}

func (c *Counters) Inc(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.values[name]++
}

func (c *Counters) Get(name string) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.values[name]
}
