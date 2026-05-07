package observability

import "testing"

func TestLogEventValid(t *testing.T) {
	event := LogEvent{Service: "backend", Action: "refresh", Status: "succeeded", RequestID: "req-1"}
	if !event.Valid() {
		t.Fatal("expected event to be valid")
	}
}

func TestCounters(t *testing.T) {
	counters := NewCounters()
	counters.Inc("refresh.succeeded")
	counters.Inc("refresh.succeeded")
	if got := counters.Get("refresh.succeeded"); got != 2 {
		t.Fatalf("unexpected counter %d", got)
	}
}
