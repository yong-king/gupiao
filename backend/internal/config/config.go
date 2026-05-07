package config

import (
	"encoding/json"
	"os"
	"time"
)

type Config struct {
	Addr         string        `json:"addr"`
	DatabaseURL  string        `json:"database_url"`
	RedisAddr    string        `json:"redis_addr"`
	AgentURL     string        `json:"agent_url"`
	Auth         AuthConfig    `json:"auth"`
	LLM          LLMConfig     `json:"llm"`
	Repository   RepoConfig    `json:"repository"`
	StockSources []StockSource `json:"stock_sources"`
}

type AuthConfig struct {
	Required          bool `json:"required"`
	SessionTTLMinutes int  `json:"session_ttl_minutes"`
}

type LLMConfig struct {
	Provider  string `json:"provider"`
	Model     string `json:"model"`
	APIKeyEnv string `json:"api_key_env"`
}

type RepoConfig struct {
	RemoteURL       string `json:"remote_url"`
	Branch          string `json:"branch"`
	AskBeforeCommit bool   `json:"ask_before_commit"`
}

type StockSource struct {
	Name               string `json:"name"`
	Type               string `json:"type"`
	BaseURL            string `json:"base_url"`
	RateLimitPerMinute int    `json:"rate_limit_per_minute"`
}

func Default() Config {
	return Config{
		Addr:        ":8080",
		DatabaseURL: "postgres://jijin:jijin@localhost:5432/jijin?sslmode=disable",
		RedisAddr:   "localhost:6379",
		AgentURL:    "http://127.0.0.1:8090",
		Auth: AuthConfig{
			Required:          true,
			SessionTTLMinutes: int((24 * time.Hour).Minutes()),
		},
		LLM: LLMConfig{
			Provider:  "deepseek",
			Model:     "deepseek-chat",
			APIKeyEnv: "DEEPSEEK_API_KEY",
		},
		Repository: RepoConfig{
			RemoteURL:       "",
			Branch:          "main",
			AskBeforeCommit: true,
		},
		StockSources: []StockSource{{Name: "stooq", Type: "stooq", BaseURL: "https://stooq.com/q/l/", RateLimitPerMinute: 20}},
	}
}

func Load() Config {
	cfg := Default()
	path := os.Getenv("JIJIN_BACKEND_CONFIG")
	if path != "" {
		if loaded, err := LoadFile(path); err == nil {
			cfg = loaded
		}
	}
	applyEnv(&cfg)
	return cfg
}

func LoadFile(path string) (Config, error) {
	cfg := Default()
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	applyEnv(&cfg)
	return cfg, nil
}

func applyEnv(cfg *Config) {
	if value := os.Getenv("GO_BACKEND_ADDR"); value != "" {
		cfg.Addr = value
	}
	if value := os.Getenv("DATABASE_URL"); value != "" {
		cfg.DatabaseURL = value
	}
	if value := os.Getenv("REDIS_ADDR"); value != "" {
		cfg.RedisAddr = value
	}
	if value := os.Getenv("AGENT_URL"); value != "" {
		cfg.AgentURL = value
	}
	if value := os.Getenv("REPOSITORY_REMOTE_URL"); value != "" {
		cfg.Repository.RemoteURL = value
	}
	if value := os.Getenv("REPOSITORY_BRANCH"); value != "" {
		cfg.Repository.Branch = value
	}
}
