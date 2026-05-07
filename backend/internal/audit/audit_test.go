package audit

import "testing"

func TestMemoryRepositoryAppendAndList(t *testing.T) {
	repo := NewMemoryRepository()
	entry := Entry{
		ID:        "audit-1",
		Action:    "settings.updated",
		Target:    "user_settings",
		TargetID:  "user-1",
		RequestID: "req-1",
		Source:    "api",
		Metadata: map[string]string{
			"refresh_mode": "conservative",
		},
	}

	if err := repo.Append(entry); err != nil {
		t.Fatalf("append audit entry: %v", err)
	}

	entries := repo.List()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Action != "settings.updated" {
		t.Fatalf("unexpected action %q", entries[0].Action)
	}
	if entries[0].CreatedAt.IsZero() {
		t.Fatal("expected created at to be set")
	}
}

func TestMemoryRepositoryCopiesMetadata(t *testing.T) {
	repo := NewMemoryRepository()
	metadata := map[string]string{"key": "original"}

	if err := repo.Append(Entry{
		ID:        "audit-1",
		Action:    "settings.updated",
		Target:    "user_settings",
		RequestID: "req-1",
		Source:    "api",
		Metadata:  metadata,
	}); err != nil {
		t.Fatalf("append audit entry: %v", err)
	}

	metadata["key"] = "mutated"
	entries := repo.List()
	entries[0].Metadata["key"] = "also-mutated"

	again := repo.List()
	if again[0].Metadata["key"] != "original" {
		t.Fatalf("metadata was not protected, got %q", again[0].Metadata["key"])
	}
}
