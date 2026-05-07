package settings

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type RefreshMode string

const (
	RefreshModeManual       RefreshMode = "manual"
	RefreshModeConservative RefreshMode = "conservative"
	RefreshModeStandard     RefreshMode = "standard"
)

type Settings struct {
	UserID                    string
	RefreshMode               RefreshMode
	DefaultCooldown           time.Duration
	EmailNotificationsEnabled bool
}

func Default(userID string) Settings {
	return Settings{
		UserID:                    userID,
		RefreshMode:               RefreshModeConservative,
		DefaultCooldown:           30 * time.Minute,
		EmailNotificationsEnabled: false,
	}
}

func (s Settings) Validate() error {
	if s.UserID == "" {
		return errors.New("user id is required")
	}
	if s.DefaultCooldown < 0 {
		return errors.New("default cooldown must not be negative")
	}

	switch s.RefreshMode {
	case RefreshModeManual, RefreshModeConservative, RefreshModeStandard:
		return nil
	default:
		return fmt.Errorf("unsupported refresh mode %q", s.RefreshMode)
	}
}

type Repository interface {
	Save(settings Settings) error
	FindByUserID(userID string) (Settings, bool)
}

type MemoryRepository struct {
	mu       sync.RWMutex
	settings map[string]Settings
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		settings: make(map[string]Settings),
	}
}

func (r *MemoryRepository) Save(settings Settings) error {
	if err := settings.Validate(); err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.settings[settings.UserID] = settings
	return nil
}

func (r *MemoryRepository) FindByUserID(userID string) (Settings, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	settings, ok := r.settings[userID]
	return settings, ok
}
