package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	data := []byte(`{"addr":":9000","database_url":"db","redis_addr":"redis","agent_url":"agent","auth":{"required":true,"session_ttl_minutes":60},"llm":{"provider":"deepseek","model":"test","api_key_env":"DEEPSEEK_API_KEY"},"repository":{"remote_url":"git@example.com:demo/jijin.git","branch":"main","ask_before_commit":true},"stock_sources":[{"name":"mock","type":"mock","rate_limit_per_minute":10}]}`)
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadFile(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if cfg.Addr != ":9000" || cfg.LLM.Model != "test" || cfg.Repository.Branch != "main" || len(cfg.StockSources) != 1 {
		t.Fatalf("unexpected config: %#v", cfg)
	}
}

func TestEnvOverrides(t *testing.T) {
	t.Setenv("DATABASE_URL", "from-env")
	cfg := Default()
	applyEnv(&cfg)
	if cfg.DatabaseURL != "from-env" {
		t.Fatalf("expected env override, got %q", cfg.DatabaseURL)
	}
}

func TestRepositoryEnvOverrides(t *testing.T) {
	t.Setenv("REPOSITORY_REMOTE_URL", "git@example.com:team/jijin.git")
	t.Setenv("REPOSITORY_BRANCH", "develop")
	cfg := Default()
	applyEnv(&cfg)
	if cfg.Repository.RemoteURL != "git@example.com:team/jijin.git" || cfg.Repository.Branch != "develop" {
		t.Fatalf("expected repository env override, got %#v", cfg.Repository)
	}
}
