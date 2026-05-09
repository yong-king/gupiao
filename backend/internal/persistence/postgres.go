package persistence

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"jijin/backend/internal/accounts"
	"jijin/backend/internal/alerts"
	"jijin/backend/internal/auth"
	"jijin/backend/internal/holdings"
	"jijin/backend/internal/marketdata"
	"jijin/backend/internal/notifications"
	"jijin/backend/internal/settings"
	"jijin/backend/internal/watchlist"
)

type Store struct {
	db *sql.DB
}

type RAGDocument struct {
	ID         string
	UserID     string
	Market     string
	Symbol     string
	SourceType string
	SourceID   string
	Content    string
	Metadata   map[string]string
	Embedding  []float64
	CreatedAt  time.Time
}

type WorkflowJob struct {
	ID             string         `json:"id"`
	UserID         string         `json:"user_id"`
	WorkflowType   string         `json:"workflow_type"`
	AttentionLevel string         `json:"attention_level"`
	Market         string         `json:"market"`
	Symbol         string         `json:"symbol"`
	Status         string         `json:"status"`
	TargetCount    int            `json:"target_count"`
	Summary        string         `json:"summary"`
	Error          string         `json:"error"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	Steps          []WorkflowStep `json:"steps"`
}

type WorkflowStep struct {
	ID            string    `json:"id"`
	JobID         string    `json:"job_id"`
	StepName      string    `json:"step_name"`
	AgentName     string    `json:"agent_name"`
	Status        string    `json:"status"`
	InputSummary  string    `json:"input_summary"`
	OutputSummary string    `json:"output_summary"`
	Error         string    `json:"error"`
	StartedAt     time.Time `json:"started_at"`
	CompletedAt   time.Time `json:"completed_at"`
}

type AssistantMessage struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	SessionID      string    `json:"session_id"`
	Market         string    `json:"market"`
	Symbol         string    `json:"symbol"`
	Question       string    `json:"question"`
	Answer         string    `json:"answer"`
	ContextSummary string    `json:"context_summary"`
	RAGDocumentIDs []string  `json:"rag_document_ids"`
	CreatedAt      time.Time `json:"created_at"`
}

type OperationLog struct {
	ID            string            `json:"id"`
	UserID        string            `json:"user_id"`
	Market        string            `json:"market"`
	Symbol        string            `json:"symbol"`
	OperationType string            `json:"operation_type"`
	Component     string            `json:"component"`
	Model         string            `json:"model"`
	InputSummary  string            `json:"input_summary"`
	OutputSummary string            `json:"output_summary"`
	Metadata      map[string]string `json:"metadata"`
	CreatedAt     time.Time         `json:"created_at"`
}

func Open(ctx context.Context, databaseURL string) (*Store, error) {
	if strings.TrimSpace(databaseURL) == "" {
		return nil, errors.New("database url is required")
	}
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Store) SaveUser(user auth.User) error {
	now := time.Now().UTC()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.Exec(`INSERT INTO users(id, display_name, created_at, updated_at) VALUES($1,$2,$3,$4)
		ON CONFLICT (id) DO UPDATE SET display_name=EXCLUDED.display_name, updated_at=EXCLUDED.updated_at`,
		user.ID, user.Email, user.CreatedAt, now); err != nil {
		return err
	}
	if _, err := tx.Exec(`INSERT INTO auth_users(id, email, password_hash, created_at) VALUES($1,$2,$3,$4)
		ON CONFLICT (id) DO UPDATE SET email=EXCLUDED.email, password_hash=EXCLUDED.password_hash`,
		user.ID, strings.ToLower(strings.TrimSpace(user.Email)), user.PasswordHash, user.CreatedAt); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Store) FindUserByEmail(email string) (auth.User, bool) {
	return s.findUser(`SELECT id, email, password_hash, created_at FROM auth_users WHERE email=$1`, strings.ToLower(strings.TrimSpace(email)))
}

func (s *Store) FindUserByID(id string) (auth.User, bool) {
	return s.findUser(`SELECT id, email, password_hash, created_at FROM auth_users WHERE id=$1`, id)
}

func (s *Store) findUser(query string, arg string) (auth.User, bool) {
	var user auth.User
	err := s.db.QueryRow(query, arg).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	return user, err == nil
}

func (s *Store) SaveSession(session auth.Session) error {
	_, err := s.db.Exec(`INSERT INTO auth_sessions(token_hash, user_id, expires_at, created_at) VALUES($1,$2,$3,$4)
		ON CONFLICT (token_hash) DO UPDATE SET user_id=EXCLUDED.user_id, expires_at=EXCLUDED.expires_at`,
		session.TokenHash, session.UserID, session.ExpiresAt, time.Now().UTC())
	return err
}

func (s *Store) FindSessionByTokenHash(hash string) (auth.Session, bool) {
	var session auth.Session
	err := s.db.QueryRow(`SELECT user_id, token_hash, expires_at FROM auth_sessions WHERE token_hash=$1`, hash).
		Scan(&session.UserID, &session.TokenHash, &session.ExpiresAt)
	return session, err == nil
}

func (s *Store) SaveWatchlist(w watchlist.Watchlist) error {
	now := time.Now().UTC()
	if w.CreatedAt.IsZero() {
		w.CreatedAt = now
	}
	_, err := s.db.Exec(`INSERT INTO watchlists(id,user_id,name,created_at,updated_at) VALUES($1,$2,$3,$4,$5)
		ON CONFLICT (id) DO UPDATE SET user_id=EXCLUDED.user_id,name=EXCLUDED.name,updated_at=EXCLUDED.updated_at`,
		w.ID, w.UserID, w.Name, w.CreatedAt, now)
	return err
}

func (s *Store) ListWatchlistsByUser(userID string) ([]watchlist.Watchlist, error) {
	rows, err := s.db.Query(`SELECT id,user_id,name,created_at,updated_at FROM watchlists WHERE user_id=$1 ORDER BY updated_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []watchlist.Watchlist{}
	for rows.Next() {
		var w watchlist.Watchlist
		if err := rows.Scan(&w.ID, &w.UserID, &w.Name, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		symbols, err := s.ListWatchlistSymbols(w.ID)
		if err != nil {
			return nil, err
		}
		w.Symbols = symbols
		out = append(out, w)
	}
	return out, rows.Err()
}

func (s *Store) FindWatchlistByID(id string) (watchlist.Watchlist, bool) {
	var w watchlist.Watchlist
	err := s.db.QueryRow(`SELECT id,user_id,name,created_at,updated_at FROM watchlists WHERE id=$1`, id).
		Scan(&w.ID, &w.UserID, &w.Name, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		return watchlist.Watchlist{}, false
	}
	symbols, err := s.ListWatchlistSymbols(id)
	if err == nil {
		w.Symbols = symbols
	}
	return w, true
}

func (s *Store) UpsertWatchlistSymbol(watchlistID string, symbol watchlist.Symbol) error {
	symbol = symbol.Normalized()
	id := watchlistID + ":" + symbol.Market + ":" + symbol.Symbol
	now := time.Now().UTC()
	_, err := s.db.Exec(`INSERT INTO watchlist_symbols(id,watchlist_id,market,symbol,note,buy_price,sell_price,created_at,updated_at)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT (watchlist_id, market, symbol) DO UPDATE SET note=EXCLUDED.note,buy_price=EXCLUDED.buy_price,sell_price=EXCLUDED.sell_price,updated_at=EXCLUDED.updated_at`,
		id, watchlistID, symbol.Market, symbol.Symbol, symbol.Note, symbol.BuyPrice, symbol.SellPrice, now, now)
	return err
}

func (s *Store) ListWatchlistSymbols(watchlistID string) ([]watchlist.Symbol, error) {
	rows, err := s.db.Query(`SELECT market,symbol,note,buy_price,sell_price FROM watchlist_symbols WHERE watchlist_id=$1 ORDER BY created_at`, watchlistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []watchlist.Symbol{}
	for rows.Next() {
		var symbol watchlist.Symbol
		if err := rows.Scan(&symbol.Market, &symbol.Symbol, &symbol.Note, &symbol.BuyPrice, &symbol.SellPrice); err != nil {
			return nil, err
		}
		out = append(out, symbol)
	}
	return out, rows.Err()
}

func (s *Store) DeleteWatchlist(userID string, id string) (bool, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback()
	if _, err := tx.Exec(`DELETE FROM watchlist_symbols WHERE watchlist_id=$1`, id); err != nil {
		return false, err
	}
	result, err := tx.Exec(`DELETE FROM watchlists WHERE id=$1 AND user_id=$2`, id, userID)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, tx.Commit()
}

func (s *Store) DeleteWatchlistSymbol(watchlistID string, market string, symbol string) (bool, error) {
	target := watchlist.Symbol{Market: market, Symbol: symbol}.Normalized()
	result, err := s.db.Exec(`DELETE FROM watchlist_symbols WHERE watchlist_id=$1 AND market=$2 AND symbol=$3`, watchlistID, target.Market, target.Symbol)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func (s *Store) UpsertHolding(item holdings.Holding) (holdings.Holding, error) {
	item = item.Normalized()
	now := time.Now().UTC()
	if item.CreatedAt.IsZero() {
		item.CreatedAt = now
	}
	id := item.UserID + ":" + item.Market + ":" + item.Symbol
	_, err := s.db.Exec(`INSERT INTO holdings(id,user_id,market,symbol,quantity,cost_basis,attention_level,created_at,updated_at)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT (id) DO UPDATE SET quantity=EXCLUDED.quantity,cost_basis=EXCLUDED.cost_basis,attention_level=EXCLUDED.attention_level,updated_at=EXCLUDED.updated_at`,
		id, item.UserID, item.Market, item.Symbol, item.Quantity, item.CostBasis, item.AttentionLevel, item.CreatedAt, now)
	return item, err
}

func (s *Store) ReplaceHoldingsForUser(userID string, items []holdings.Holding) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.Exec(`DELETE FROM holdings WHERE user_id=$1`, userID); err != nil {
		return err
	}
	now := time.Now().UTC()
	for _, item := range items {
		item = item.Normalized()
		id := userID + ":" + item.Market + ":" + item.Symbol
		if _, err := tx.Exec(`INSERT INTO holdings(id,user_id,market,symbol,quantity,cost_basis,attention_level,created_at,updated_at) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
			id, userID, item.Market, item.Symbol, item.Quantity, item.CostBasis, item.AttentionLevel, now, now); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) ListHoldingsByUser(userID string) ([]holdings.Holding, error) {
	rows, err := s.db.Query(`SELECT id,user_id,market,symbol,quantity,cost_basis,attention_level,created_at,updated_at FROM holdings WHERE user_id=$1 ORDER BY updated_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []holdings.Holding{}
	for rows.Next() {
		var item holdings.Holding
		if err := rows.Scan(&item.ID, &item.UserID, &item.Market, &item.Symbol, &item.Quantity, &item.CostBasis, &item.AttentionLevel, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (s *Store) DeleteHolding(userID string, market string, symbol string) (bool, error) {
	normalized := holdings.Holding{UserID: userID, Market: market, Symbol: symbol}.Normalized()
	id := normalized.UserID + ":" + normalized.Market + ":" + normalized.Symbol
	result, err := s.db.Exec(`DELETE FROM holdings WHERE id=$1`, id)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func (s *Store) SaveAlertRule(rule alerts.Rule) error {
	now := time.Now().UTC()
	if rule.CreatedAt.IsZero() {
		rule.CreatedAt = now
	}
	_, err := s.db.Exec(`INSERT INTO alert_rules(id,user_id,market,symbol,type,threshold,signal,risk_level,enabled,cooldown_seconds,created_at,updated_at)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		ON CONFLICT (id) DO UPDATE SET type=EXCLUDED.type,threshold=EXCLUDED.threshold,signal=EXCLUDED.signal,risk_level=EXCLUDED.risk_level,enabled=EXCLUDED.enabled,cooldown_seconds=EXCLUDED.cooldown_seconds,updated_at=EXCLUDED.updated_at`,
		rule.ID, rule.UserID, rule.Market, rule.Symbol, string(rule.Type), rule.Threshold, string(rule.Signal), string(rule.RiskLevel), rule.Enabled, int(rule.Cooldown.Seconds()), rule.CreatedAt, now)
	return err
}

func (s *Store) ListAlertRulesByUser(userID string) ([]alerts.Rule, error) {
	rows, err := s.db.Query(`SELECT id,user_id,market,symbol,type,threshold,signal,risk_level,enabled,cooldown_seconds,created_at,updated_at FROM alert_rules WHERE user_id=$1 ORDER BY updated_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []alerts.Rule{}
	for rows.Next() {
		var rule alerts.Rule
		var ruleType, signal, risk string
		var cooldownSeconds int
		if err := rows.Scan(&rule.ID, &rule.UserID, &rule.Market, &rule.Symbol, &ruleType, &rule.Threshold, &signal, &risk, &rule.Enabled, &cooldownSeconds, &rule.CreatedAt, &rule.UpdatedAt); err != nil {
			return nil, err
		}
		rule.Type = alerts.RuleType(ruleType)
		rule.Signal = alerts.Signal(signal)
		rule.RiskLevel = alerts.RiskLevel(risk)
		rule.Cooldown = time.Duration(cooldownSeconds) * time.Second
		out = append(out, rule)
	}
	return out, rows.Err()
}

func (s *Store) DeleteAlertRule(userID string, id string) (bool, error) {
	result, err := s.db.Exec(`DELETE FROM alert_rules WHERE id=$1 AND user_id=$2`, id, userID)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	return affected > 0, nil
}

func (s *Store) SaveAlertEvent(event alerts.Event) error {
	_, err := s.db.Exec(`INSERT INTO alert_events(id,user_id,rule_id,market,symbol,signal,risk_level,summary,source,data_time,read,created_at)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		ON CONFLICT (id) DO NOTHING`,
		event.ID, event.UserID, event.RuleID, event.Market, event.Symbol, string(event.Signal), string(event.RiskLevel), event.Summary, event.Source, event.DataTime, event.Read, time.Now().UTC())
	return err
}

func (s *Store) ListAlertEventsByUser(userID string) ([]alerts.Event, error) {
	rows, err := s.db.Query(`SELECT id,user_id,rule_id,market,symbol,signal,risk_level,summary,source,data_time,read,created_at FROM alert_events WHERE user_id=$1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []alerts.Event{}
	for rows.Next() {
		var event alerts.Event
		var signal, risk string
		if err := rows.Scan(&event.ID, &event.UserID, &event.RuleID, &event.Market, &event.Symbol, &signal, &risk, &event.Summary, &event.Source, &event.DataTime, &event.Read, &event.CreatedAt); err != nil {
			return nil, err
		}
		event.Signal = alerts.Signal(signal)
		event.RiskLevel = alerts.RiskLevel(risk)
		out = append(out, event)
	}
	return out, rows.Err()
}

func (s *Store) SaveNotification(message notifications.Message) error {
	_, err := s.db.Exec(`INSERT INTO notifications(id,user_id,title,summary,signal,risk_level,market,symbol,data_time,read,created_at)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		ON CONFLICT (id) DO UPDATE SET read=EXCLUDED.read`,
		message.ID, message.UserID, message.Title, message.Summary, string(message.Signal), string(message.RiskLevel), message.Market, message.Symbol, message.DataTime, message.Read, time.Now().UTC())
	return err
}

func (s *Store) ListNotificationsByUser(userID string) ([]notifications.Message, error) {
	rows, err := s.db.Query(`SELECT id,user_id,title,summary,signal,risk_level,market,symbol,data_time,read,created_at FROM notifications WHERE user_id=$1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []notifications.Message{}
	for rows.Next() {
		var message notifications.Message
		var signal, risk string
		if err := rows.Scan(&message.ID, &message.UserID, &message.Title, &message.Summary, &signal, &risk, &message.Market, &message.Symbol, &message.DataTime, &message.Read, &message.CreatedAt); err != nil {
			return nil, err
		}
		message.Signal = alerts.Signal(signal)
		message.RiskLevel = alerts.RiskLevel(risk)
		out = append(out, message)
	}
	return out, rows.Err()
}

func (s *Store) MarkNotificationRead(id string) error {
	result, err := s.db.Exec(`UPDATE notifications SET read=TRUE WHERE id=$1`, id)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return errors.New("message not found")
	}
	return nil
}

func (s *Store) SaveAccount(config accounts.Config) error {
	if config.Metadata == nil {
		config.Metadata = map[string]string{}
	}
	metadata, err := json.Marshal(config.Metadata)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	if config.CreatedAt.IsZero() {
		config.CreatedAt = now
	}
	_, err = s.db.Exec(`INSERT INTO broker_accounts(id,user_id,alias,refresh_mode,read_only,metadata,created_at,updated_at)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8)
		ON CONFLICT (id) DO UPDATE SET alias=EXCLUDED.alias,refresh_mode=EXCLUDED.refresh_mode,read_only=EXCLUDED.read_only,metadata=EXCLUDED.metadata,updated_at=EXCLUDED.updated_at`,
		config.ID, config.UserID, config.Alias, string(config.RefreshMode), config.ReadOnly, string(metadata), config.CreatedAt, now)
	return err
}

func (s *Store) ListAccountsByUser(userID string) ([]accounts.Config, error) {
	rows, err := s.db.Query(`SELECT id,user_id,alias,refresh_mode,read_only,metadata,created_at,updated_at FROM broker_accounts WHERE user_id=$1 ORDER BY updated_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []accounts.Config{}
	for rows.Next() {
		var config accounts.Config
		var mode string
		var metadata string
		if err := rows.Scan(&config.ID, &config.UserID, &config.Alias, &mode, &config.ReadOnly, &metadata, &config.CreatedAt, &config.UpdatedAt); err != nil {
			return nil, err
		}
		config.RefreshMode = settings.RefreshMode(mode)
		_ = json.Unmarshal([]byte(metadata), &config.Metadata)
		out = append(out, config)
	}
	return out, rows.Err()
}

func (s *Store) SaveSnapshots(snapshots []marketdata.Snapshot) error {
	for _, snapshot := range snapshots {
		id := snapshot.Market + ":" + snapshot.Symbol + ":" + snapshot.DataTime.Format(time.RFC3339Nano)
		if _, err := s.db.Exec(`INSERT INTO price_snapshots(id,market,symbol,name,open,high,low,price,previous_close,change_percent,volume,source,data_time,created_at)
			VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)
			ON CONFLICT (id) DO UPDATE SET price=EXCLUDED.price,previous_close=EXCLUDED.previous_close,change_percent=EXCLUDED.change_percent,volume=EXCLUDED.volume`,
			id, snapshot.Market, snapshot.Symbol, snapshot.Name, snapshot.Open, snapshot.High, snapshot.Low, snapshot.Price, snapshot.PreviousClose, snapshot.ChangePercent, snapshot.Volume, snapshot.Source, snapshot.DataTime, time.Now().UTC()); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) ListSnapshots(market string, symbol string) ([]marketdata.Snapshot, error) {
	rows, err := s.db.Query(`SELECT market,symbol,name,open,high,low,price,previous_close,change_percent,volume,source,data_time,created_at FROM price_snapshots WHERE market=$1 AND symbol=$2 ORDER BY data_time`, market, symbol)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []marketdata.Snapshot{}
	for rows.Next() {
		var snapshot marketdata.Snapshot
		if err := rows.Scan(&snapshot.Market, &snapshot.Symbol, &snapshot.Name, &snapshot.Open, &snapshot.High, &snapshot.Low, &snapshot.Price, &snapshot.PreviousClose, &snapshot.ChangePercent, &snapshot.Volume, &snapshot.Source, &snapshot.DataTime, &snapshot.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, snapshot)
	}
	return out, rows.Err()
}

func (s *Store) SaveRAGDocument(document RAGDocument) error {
	if document.CreatedAt.IsZero() {
		document.CreatedAt = time.Now().UTC()
	}
	if document.Metadata == nil {
		document.Metadata = map[string]string{}
	}
	if len(document.Embedding) == 0 {
		document.Embedding = []float64{}
	}
	metadata, err := json.Marshal(document.Metadata)
	if err != nil {
		return err
	}
	embedding, err := json.Marshal(document.Embedding)
	if err != nil {
		return err
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.Exec(`INSERT INTO rag_documents(id,user_id,market,symbol,source_type,source_id,content,metadata,created_at)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)
		ON CONFLICT (id) DO UPDATE SET content=EXCLUDED.content,metadata=EXCLUDED.metadata`,
		document.ID, document.UserID, document.Market, document.Symbol, document.SourceType, document.SourceID, document.Content, string(metadata), document.CreatedAt); err != nil {
		return err
	}
	vectorID := document.ID + ":vector"
	if _, err := tx.Exec(`INSERT INTO rag_vectors(id,rag_document_id,provider,model,embedding,status,created_at)
		VALUES($1,$2,$3,$4,$5,$6,$7)
		ON CONFLICT (id) DO UPDATE SET embedding=EXCLUDED.embedding,status=EXCLUDED.status`,
		vectorID, document.ID, "local", "deterministic-hash-v1", string(embedding), "indexed", document.CreatedAt); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Store) ListRAGDocuments(userID string, market string, symbol string, limit int) ([]RAGDocument, error) {
	if limit <= 0 {
		limit = 5
	}
	rows, err := s.db.Query(`SELECT id,user_id,market,symbol,source_type,source_id,content,metadata,created_at
		FROM rag_documents
		WHERE user_id=$1 AND market=$2 AND symbol=$3
		ORDER BY created_at DESC
		LIMIT $4`, userID, market, symbol, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []RAGDocument{}
	for rows.Next() {
		var document RAGDocument
		var metadata string
		if err := rows.Scan(&document.ID, &document.UserID, &document.Market, &document.Symbol, &document.SourceType, &document.SourceID, &document.Content, &metadata, &document.CreatedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal([]byte(metadata), &document.Metadata)
		out = append(out, document)
	}
	return out, rows.Err()
}

func (s *Store) SaveWorkflowJob(job WorkflowJob) error {
	now := time.Now().UTC()
	if job.CreatedAt.IsZero() {
		job.CreatedAt = now
	}
	if job.UpdatedAt.IsZero() {
		job.UpdatedAt = now
	}
	_, err := s.db.Exec(`INSERT INTO agent_workflow_jobs(id,user_id,workflow_type,attention_level,market,symbol,status,target_count,summary,error,created_at,updated_at)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		ON CONFLICT (id) DO UPDATE SET status=EXCLUDED.status,target_count=EXCLUDED.target_count,summary=EXCLUDED.summary,error=EXCLUDED.error,updated_at=EXCLUDED.updated_at`,
		job.ID, job.UserID, job.WorkflowType, job.AttentionLevel, job.Market, job.Symbol, job.Status, job.TargetCount, job.Summary, job.Error, job.CreatedAt, job.UpdatedAt)
	return err
}

func (s *Store) SaveWorkflowStep(step WorkflowStep) error {
	now := time.Now().UTC()
	if step.StartedAt.IsZero() {
		step.StartedAt = now
	}
	if step.CompletedAt.IsZero() {
		step.CompletedAt = now
	}
	_, err := s.db.Exec(`INSERT INTO agent_workflow_steps(id,job_id,step_name,agent_name,status,input_summary,output_summary,error,started_at,completed_at)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		ON CONFLICT (id) DO UPDATE SET status=EXCLUDED.status,output_summary=EXCLUDED.output_summary,error=EXCLUDED.error,completed_at=EXCLUDED.completed_at`,
		step.ID, step.JobID, step.StepName, step.AgentName, step.Status, step.InputSummary, step.OutputSummary, step.Error, step.StartedAt, step.CompletedAt)
	return err
}

func (s *Store) ListWorkflowJobs(userID string, limit int) ([]WorkflowJob, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := s.db.Query(`SELECT id,user_id,workflow_type,attention_level,market,symbol,status,target_count,summary,error,created_at,updated_at
		FROM agent_workflow_jobs WHERE user_id=$1 ORDER BY created_at DESC LIMIT $2`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []WorkflowJob{}
	for rows.Next() {
		var job WorkflowJob
		if err := rows.Scan(&job.ID, &job.UserID, &job.WorkflowType, &job.AttentionLevel, &job.Market, &job.Symbol, &job.Status, &job.TargetCount, &job.Summary, &job.Error, &job.CreatedAt, &job.UpdatedAt); err != nil {
			return nil, err
		}
		steps, err := s.ListWorkflowSteps(job.ID)
		if err != nil {
			return nil, err
		}
		job.Steps = steps
		out = append(out, job)
	}
	return out, rows.Err()
}

func (s *Store) ListWorkflowSteps(jobID string) ([]WorkflowStep, error) {
	rows, err := s.db.Query(`SELECT id,job_id,step_name,agent_name,status,input_summary,output_summary,error,started_at,completed_at
		FROM agent_workflow_steps WHERE job_id=$1 ORDER BY started_at`, jobID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []WorkflowStep{}
	for rows.Next() {
		var step WorkflowStep
		if err := rows.Scan(&step.ID, &step.JobID, &step.StepName, &step.AgentName, &step.Status, &step.InputSummary, &step.OutputSummary, &step.Error, &step.StartedAt, &step.CompletedAt); err != nil {
			return nil, err
		}
		out = append(out, step)
	}
	return out, rows.Err()
}

func (s *Store) SaveAssistantMessage(message AssistantMessage) error {
	if message.CreatedAt.IsZero() {
		message.CreatedAt = time.Now().UTC()
	}
	if strings.TrimSpace(message.SessionID) == "" {
		message.SessionID = "default"
	}
	ragIDs, err := json.Marshal(message.RAGDocumentIDs)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`INSERT INTO stock_assistant_messages(id,user_id,session_id,market,symbol,question,answer,context_summary,rag_document_ids,created_at)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		message.ID, message.UserID, message.SessionID, message.Market, message.Symbol, message.Question, message.Answer, message.ContextSummary, string(ragIDs), message.CreatedAt)
	return err
}

func (s *Store) ListAssistantMessages(userID string, sessionID string, market string, symbol string, limit int) ([]AssistantMessage, error) {
	if limit <= 0 {
		limit = 8
	}
	if strings.TrimSpace(sessionID) == "" {
		sessionID = "default"
	}
	rows, err := s.db.Query(`SELECT id,user_id,session_id,market,symbol,question,answer,context_summary,rag_document_ids,created_at
		FROM stock_assistant_messages
		WHERE user_id=$1 AND session_id=$2 AND market=$3 AND symbol=$4
		ORDER BY created_at DESC
		LIMIT $5`, userID, sessionID, market, symbol, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []AssistantMessage{}
	for rows.Next() {
		var message AssistantMessage
		var ragIDs string
		if err := rows.Scan(&message.ID, &message.UserID, &message.SessionID, &message.Market, &message.Symbol, &message.Question, &message.Answer, &message.ContextSummary, &ragIDs, &message.CreatedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal([]byte(ragIDs), &message.RAGDocumentIDs)
		out = append(out, message)
	}
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out, rows.Err()
}

func (s *Store) SaveOperationLog(log OperationLog) error {
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now().UTC()
	}
	if log.Metadata == nil {
		log.Metadata = map[string]string{}
	}
	metadata, err := json.Marshal(log.Metadata)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`INSERT INTO operation_logs(id,user_id,market,symbol,operation_type,component,model,input_summary,output_summary,metadata,created_at)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		ON CONFLICT (id) DO UPDATE SET output_summary=EXCLUDED.output_summary,metadata=EXCLUDED.metadata`,
		log.ID, log.UserID, log.Market, log.Symbol, log.OperationType, log.Component, log.Model, log.InputSummary, log.OutputSummary, string(metadata), log.CreatedAt)
	return err
}

func (s *Store) ListOperationLogs(userID string, limit int) ([]OperationLog, error) {
	if limit <= 0 {
		limit = 100
	}
	rows, err := s.db.Query(`SELECT id,user_id,market,symbol,operation_type,component,model,input_summary,output_summary,metadata,created_at
		FROM operation_logs
		WHERE user_id=$1 OR user_id=''
		ORDER BY created_at DESC
		LIMIT $2`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []OperationLog{}
	for rows.Next() {
		var log OperationLog
		var metadata string
		if err := rows.Scan(&log.ID, &log.UserID, &log.Market, &log.Symbol, &log.OperationType, &log.Component, &log.Model, &log.InputSummary, &log.OutputSummary, &metadata, &log.CreatedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal([]byte(metadata), &log.Metadata)
		out = append(out, log)
	}
	return out, rows.Err()
}
