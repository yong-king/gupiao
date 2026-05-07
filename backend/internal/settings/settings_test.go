package settings

import "testing"

func TestDefaultSettingsUseConservativeRefresh(t *testing.T) {
	got := Default("user-1")

	if got.RefreshMode != RefreshModeConservative {
		t.Fatalf("expected conservative refresh, got %q", got.RefreshMode)
	}
	if err := got.Validate(); err != nil {
		t.Fatalf("default settings should validate: %v", err)
	}
}

func TestValidateRejectsInvalidRefreshMode(t *testing.T) {
	got := Default("user-1")
	got.RefreshMode = "aggressive"

	if err := got.Validate(); err == nil {
		t.Fatal("expected invalid refresh mode to fail")
	}
}

func TestMemoryRepositorySaveAndFind(t *testing.T) {
	repo := NewMemoryRepository()
	want := Default("user-1")

	if err := repo.Save(want); err != nil {
		t.Fatalf("save settings: %v", err)
	}

	got, ok := repo.FindByUserID("user-1")
	if !ok {
		t.Fatal("expected settings to be found")
	}
	if got != want {
		t.Fatalf("unexpected settings: %#v", got)
	}
}
