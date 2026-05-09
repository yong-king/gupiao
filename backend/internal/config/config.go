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
	Cadence      CadenceConfig `json:"cadence"`
	Repository   RepoConfig    `json:"repository"`
	StockSources []StockSource `json:"stock_sources"`
}

type AuthConfig struct {
	Required          bool `json:"required"`
	SessionTTLMinutes int  `json:"session_ttl_minutes"`
}

type LLMConfig struct {
	Provider   string `json:"provider"`
	Model      string `json:"model"`
	ChatModel  string `json:"chat_model"`
	FlashModel string `json:"flash_model"`
	ProModel   string `json:"pro_model"`
	APIKeyEnv  string `json:"api_key_env"`
}

type CadenceConfig struct {
	ProductResearch map[string]string `json:"product_research"`
	RealtimeQuote   map[string]string `json:"realtime_quote"`
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
			Provider:   "deepseek",
			Model:      "deepseek-chat",
			ChatModel:  "deepseek-chat",
			FlashModel: "deepseek-v4-flash",
			ProModel:   "deepseek-v4-pro",
			APIKeyEnv:  "DEEPSEEK_API_KEY",
		},
		Cadence: CadenceConfig{
			ProductResearch: map[string]string{"high": "1h", "medium": "2h", "low": "4h"},
			RealtimeQuote:   map[string]string{"high": "2m", "medium": "5m", "low": "10m"},
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

func (c Config) ProductResearchInterval(level string) time.Duration {
	return durationForAttention(c.Cadence.ProductResearch, level, map[string]string{"high": "1h", "medium": "2h", "low": "4h"})
}

func (c Config) RealtimeQuoteInterval(level string) time.Duration {
	return durationForAttention(c.Cadence.RealtimeQuote, level, map[string]string{"high": "2m", "medium": "5m", "low": "10m"})
}

func durationForAttention(values map[string]string, level string, defaults map[string]string) time.Duration {
	key := level
	if key != "high" && key != "low" {
		key = "medium"
	}
	raw := values[key]
	if raw == "" {
		raw = defaults[key]
	}
	duration, err := time.ParseDuration(raw)
	if err != nil {
		duration, _ = time.ParseDuration(defaults[key])
	}
	return duration
}
