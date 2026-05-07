package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteError(t *testing.T) {
	rec := httptest.NewRecorder()

	WriteError(rec, http.StatusBadRequest, "validation_error", "Invalid request.", "req_test")

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var body ErrorResponse
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if body.Error.Code != "validation_error" {
		t.Fatalf("unexpected error code %q", body.Error.Code)
	}
	if body.Error.RequestID != "req_test" {
		t.Fatalf("unexpected request id %q", body.Error.RequestID)
	}
}
