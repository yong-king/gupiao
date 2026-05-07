package api

import (
	"net/http"
	"os"
	"strings"

	"jijin/backend/internal/alerts"
	"jijin/backend/internal/audit"
	"jijin/backend/internal/auth"
	"jijin/backend/internal/config"
	"jijin/backend/internal/health"
	"jijin/backend/internal/holdings"
	"jijin/backend/internal/marketdata"
	"jijin/backend/internal/notifications"
	"jijin/backend/internal/refresh"
	"jijin/backend/internal/watchlist"
)

func NewRouter() http.Handler {
	server := NewServerWithConfig(watchlist.NewRepository(), holdings.NewRepository(), config.Load())
	server.authRequired = true
	return server.Routes()
}

type Server struct {
	watchlists   *watchlist.Repository
	holdings     *holdings.Repository
	audits       *audit.MemoryRepository
	refreshes    *refresh.Service
	alertRules   *alerts.RuleRepository
	alerts       *alerts.EventRepository
	notifier     *notifications.Center
	auth         *auth.Service
	cfg          config.Config
	authRequired bool
}

func NewServer(watchlists *watchlist.Repository, holdings *holdings.Repository) *Server {
	provider := marketdata.NewMockProvider()
	seedDemoQuotes(provider)
	return NewServerWithProvider(watchlists, holdings, provider)
}

func NewServerWithConfig(watchlists *watchlist.Repository, holdings *holdings.Repository, cfg config.Config) *Server {
	server := NewServerWithProvider(watchlists, holdings, marketProviderFromConfig(cfg))
	server.cfg = cfg
	return server
}

func NewServerWithProvider(watchlists *watchlist.Repository, holdings *holdings.Repository, provider marketdata.Provider) *Server {
	snapshots := marketdata.NewSnapshotRepository()
	jobs := refresh.NewJobRepository()
	return &Server{
		watchlists: watchlists,
		holdings:   holdings,
		audits:     audit.NewMemoryRepository(),
		refreshes:  refresh.NewService(provider, snapshots, jobs),
		alertRules: alerts.NewRuleRepository(),
		alerts:     alerts.NewEventRepository(),
		notifier:   notifications.NewCenter(),
		auth:       auth.NewService(),
		cfg:        config.Default(),
	}
}

func marketProviderFromConfig(cfg config.Config) marketdata.Provider {
	if len(cfg.StockSources) == 0 {
		return seededMockProvider()
	}
	source := cfg.StockSources[0]
	if strings.EqualFold(source.Type, "stooq") {
		return marketdata.NewStooqProvider(source.BaseURL)
	}
	return seededMockProvider()
}

func seededMockProvider() *marketdata.MockProvider {
	provider := marketdata.NewMockProvider()
	seedDemoQuotes(provider)
	return provider
}

func seedDemoQuotes(provider *marketdata.MockProvider) {
	provider.SetQuote(marketdata.Snapshot{Market: "US", Symbol: "AAPL", Price: 150, PreviousClose: 155, ChangePercent: -3.23, Volume: 1200000, Source: "mock"})
	provider.SetQuote(marketdata.Snapshot{Market: "US", Symbol: "MSFT", Price: 410, PreviousClose: 400, ChangePercent: 2.5, Volume: 980000, Source: "mock"})
	provider.SetQuote(marketdata.Snapshot{Market: "HK", Symbol: "0700", Price: 300, PreviousClose: 295, ChangePercent: 1.69, Volume: 860000, Source: "mock"})
}

func NewServerWithRefresh(watchlists *watchlist.Repository, holdings *holdings.Repository, refreshes *refresh.Service) *Server {
	return &Server{
		watchlists: watchlists,
		holdings:   holdings,
		audits:     audit.NewMemoryRepository(),
		refreshes:  refreshes,
		alertRules: alerts.NewRuleRepository(),
		alerts:     alerts.NewEventRepository(),
		notifier:   notifications.NewCenter(),
		auth:       auth.NewService(),
		cfg:        config.Default(),
	}
}

func (s *Server) AuditEntries() []audit.Entry {
	return s.audits.List()
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/healthz", health.Handler())
	mux.HandleFunc("/api/auth/register", s.handleRegister)
	mux.HandleFunc("/api/auth/login", s.handleLogin)
	mux.HandleFunc("/api/watchlists", s.requireAuth(s.handleWatchlists))
	mux.HandleFunc("/api/watchlists/", s.requireAuth(s.handleWatchlistByID))
	mux.HandleFunc("/api/holdings", s.requireAuth(s.handleHoldings))
	mux.HandleFunc("/api/holdings/import", s.requireAuth(s.handleHoldingsImport))
	mux.HandleFunc("/api/refresh/manual", s.requireAuth(s.handleManualRefresh))
	mux.HandleFunc("/api/market/collect", s.requireAuth(s.handleMarketCollect))
	mux.HandleFunc("/api/market/snapshots", s.requireAuth(s.handleMarketSnapshots))
	mux.HandleFunc("/api/market/daily-changes", s.requireAuth(s.handleDailyChanges))
	mux.HandleFunc("/api/stocks/profile", s.requireAuth(s.handleStockProfile))
	mux.HandleFunc("/api/alert-rules", s.requireAuth(s.handleAlertRules))
	mux.HandleFunc("/api/alerts", s.requireAuth(s.handleAlerts))
	mux.HandleFunc("/api/notifications", s.requireAuth(s.handleNotifications))
	mux.HandleFunc("/api/system/dependencies", s.requireAuth(s.handleSystemDependencies))
	return withCORS(mux)
}

func (s *Server) llmKeyConfigured() bool {
	return strings.TrimSpace(s.cfg.LLM.APIKeyEnv) != "" && os.Getenv(s.cfg.LLM.APIKeyEnv) != ""
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Private-Network", "true")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.authRequired {
			next(w, r)
			return
		}
		header := r.Header.Get("Authorization")
		token, ok := strings.CutPrefix(header, "Bearer ")
		if !ok || strings.TrimSpace(token) == "" {
			WriteError(w, http.StatusUnauthorized, "validation_error", "Authentication required.", requestID(r))
			return
		}
		if _, ok := s.auth.Authenticate(strings.TrimSpace(token)); !ok {
			WriteError(w, http.StatusUnauthorized, "validation_error", "Invalid or expired token.", requestID(r))
			return
		}
		next(w, r)
	}
}
