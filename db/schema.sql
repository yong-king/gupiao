CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    display_name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS auth_users (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS auth_sessions (
    token_hash TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES auth_users(id),
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS user_settings (
    user_id TEXT PRIMARY KEY REFERENCES users(id),
    refresh_mode TEXT NOT NULL,
    default_cooldown_seconds INTEGER NOT NULL,
    email_notifications_enabled BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS watchlists (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id),
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS watchlist_symbols (
    id TEXT PRIMARY KEY,
    watchlist_id TEXT NOT NULL REFERENCES watchlists(id),
    market TEXT NOT NULL,
    symbol TEXT NOT NULL,
    note TEXT NOT NULL DEFAULT '',
    buy_price NUMERIC NOT NULL DEFAULT 0,
    sell_price NUMERIC NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    UNIQUE (watchlist_id, market, symbol)
);

CREATE TABLE IF NOT EXISTS holdings (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id),
    market TEXT NOT NULL,
    symbol TEXT NOT NULL,
    quantity NUMERIC NOT NULL,
    cost_basis NUMERIC NOT NULL,
    attention_level TEXT NOT NULL DEFAULT 'medium',
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS refresh_jobs (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id),
    scope TEXT NOT NULL,
    status TEXT NOT NULL,
    error TEXT NOT NULL DEFAULT '',
    requested_at TIMESTAMPTZ NOT NULL,
    started_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS price_snapshots (
    id TEXT PRIMARY KEY,
    market TEXT NOT NULL,
    symbol TEXT NOT NULL,
    name TEXT NOT NULL DEFAULT '',
    open NUMERIC NOT NULL DEFAULT 0,
    high NUMERIC NOT NULL DEFAULT 0,
    low NUMERIC NOT NULL DEFAULT 0,
    price NUMERIC NOT NULL,
    volume BIGINT NOT NULL,
    source TEXT NOT NULL,
    data_time TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS daily_price_changes (
    id TEXT PRIMARY KEY,
    market TEXT NOT NULL,
    symbol TEXT NOT NULL,
    date DATE NOT NULL,
    open NUMERIC NOT NULL,
    close NUMERIC NOT NULL,
    previous_close NUMERIC NOT NULL,
    change NUMERIC NOT NULL,
    change_percent NUMERIC NOT NULL,
    volume BIGINT NOT NULL,
    source TEXT NOT NULL,
    rag_text TEXT NOT NULL,
    data_time TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    UNIQUE (market, symbol, date)
);

CREATE TABLE IF NOT EXISTS alert_rules (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id),
    market TEXT NOT NULL,
    symbol TEXT NOT NULL,
    type TEXT NOT NULL,
    threshold NUMERIC NOT NULL,
    signal TEXT NOT NULL,
    risk_level TEXT NOT NULL,
    enabled BOOLEAN NOT NULL,
    cooldown_seconds INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS alert_events (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id),
    rule_id TEXT NOT NULL REFERENCES alert_rules(id),
    market TEXT NOT NULL,
    symbol TEXT NOT NULL,
    signal TEXT NOT NULL,
    risk_level TEXT NOT NULL,
    summary TEXT NOT NULL,
    source TEXT NOT NULL,
    data_time TIMESTAMPTZ NOT NULL,
    read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS notifications (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id),
    title TEXT NOT NULL,
    summary TEXT NOT NULL,
    signal TEXT NOT NULL,
    risk_level TEXT NOT NULL,
    market TEXT NOT NULL,
    symbol TEXT NOT NULL,
    data_time TIMESTAMPTZ NOT NULL,
    read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS broker_accounts (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id),
    alias TEXT NOT NULL,
    refresh_mode TEXT NOT NULL,
    read_only BOOLEAN NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS agent_runs (
    id TEXT PRIMARY KEY,
    agent_version TEXT NOT NULL,
    input JSONB NOT NULL,
    output JSONB NOT NULL,
    status TEXT NOT NULL,
    error TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS analysis_reports (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id),
    report_date DATE NOT NULL,
    summary TEXT NOT NULL,
    risk_points JSONB NOT NULL DEFAULT '[]',
    needs_confirmation JSONB NOT NULL DEFAULT '[]',
    data_time TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS stock_profiles (
    id TEXT PRIMARY KEY,
    market TEXT NOT NULL,
    symbol TEXT NOT NULL,
    name TEXT NOT NULL,
    sector TEXT NOT NULL,
    products JSONB NOT NULL DEFAULT '[]'::jsonb,
    business TEXT NOT NULL,
    data_source TEXT NOT NULL,
    analysis TEXT NOT NULL,
    recommendation TEXT NOT NULL,
    disclaimer TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    UNIQUE (market, symbol)
);

CREATE TABLE IF NOT EXISTS rag_documents (
    id TEXT PRIMARY KEY,
    user_id TEXT REFERENCES users(id),
    market TEXT NOT NULL,
    symbol TEXT NOT NULL,
    source_type TEXT NOT NULL,
    source_id TEXT NOT NULL,
    content TEXT NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS llm_analysis_runs (
    id TEXT PRIMARY KEY,
    user_id TEXT REFERENCES users(id),
    market TEXT NOT NULL,
    symbol TEXT NOT NULL,
    provider TEXT NOT NULL,
    model TEXT NOT NULL,
    prompt_hash TEXT NOT NULL,
    input_summary TEXT NOT NULL,
    output_summary TEXT NOT NULL,
    status TEXT NOT NULL,
    error TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS rag_vectors (
    id TEXT PRIMARY KEY,
    rag_document_id TEXT NOT NULL REFERENCES rag_documents(id) ON DELETE CASCADE,
    provider TEXT NOT NULL,
    model TEXT NOT NULL,
    embedding JSONB NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS agent_workflow_jobs (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id),
    workflow_type TEXT NOT NULL,
    attention_level TEXT NOT NULL,
    market TEXT NOT NULL DEFAULT '',
    symbol TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL,
    target_count INTEGER NOT NULL DEFAULT 0,
    summary TEXT NOT NULL DEFAULT '',
    error TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS agent_workflow_steps (
    id TEXT PRIMARY KEY,
    job_id TEXT NOT NULL REFERENCES agent_workflow_jobs(id) ON DELETE CASCADE,
    step_name TEXT NOT NULL,
    agent_name TEXT NOT NULL,
    status TEXT NOT NULL,
    input_summary TEXT NOT NULL DEFAULT '',
    output_summary TEXT NOT NULL DEFAULT '',
    error TEXT NOT NULL DEFAULT '',
    started_at TIMESTAMPTZ NOT NULL,
    completed_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS stock_assistant_messages (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id),
    session_id TEXT NOT NULL DEFAULT 'default',
    market TEXT NOT NULL,
    symbol TEXT NOT NULL,
    question TEXT NOT NULL,
    answer TEXT NOT NULL,
    context_summary TEXT NOT NULL DEFAULT '',
    rag_document_ids JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS operation_logs (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL DEFAULT '',
    market TEXT NOT NULL DEFAULT '',
    symbol TEXT NOT NULL DEFAULT '',
    operation_type TEXT NOT NULL,
    component TEXT NOT NULL,
    model TEXT NOT NULL DEFAULT '',
    input_summary TEXT NOT NULL DEFAULT '',
    output_summary TEXT NOT NULL DEFAULT '',
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_watchlists_user_id ON watchlists(user_id);
CREATE INDEX IF NOT EXISTS idx_holdings_user_id ON holdings(user_id);
CREATE INDEX IF NOT EXISTS idx_price_snapshots_symbol_time ON price_snapshots(market, symbol, data_time DESC);
CREATE INDEX IF NOT EXISTS idx_refresh_jobs_user_id ON refresh_jobs(user_id);
CREATE INDEX IF NOT EXISTS idx_daily_price_changes_symbol_date ON daily_price_changes(market, symbol, date DESC);
CREATE INDEX IF NOT EXISTS idx_alert_rules_user_symbol ON alert_rules(user_id, market, symbol);
CREATE INDEX IF NOT EXISTS idx_alert_events_user_created ON alert_events(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_notifications_user_created ON notifications(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_analysis_reports_user_date ON analysis_reports(user_id, report_date DESC);
CREATE INDEX IF NOT EXISTS idx_broker_accounts_user_id ON broker_accounts(user_id);
CREATE INDEX IF NOT EXISTS idx_rag_documents_symbol_created ON rag_documents(market, symbol, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_llm_analysis_runs_user_symbol ON llm_analysis_runs(user_id, market, symbol, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_rag_vectors_document ON rag_vectors(rag_document_id);
CREATE INDEX IF NOT EXISTS idx_agent_workflow_jobs_user_created ON agent_workflow_jobs(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_agent_workflow_steps_job ON agent_workflow_steps(job_id, started_at);
CREATE INDEX IF NOT EXISTS idx_stock_assistant_messages_user_symbol ON stock_assistant_messages(user_id, market, symbol, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_stock_assistant_messages_session_symbol ON stock_assistant_messages(user_id, session_id, market, symbol, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_operation_logs_user_created ON operation_logs(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_operation_logs_symbol_created ON operation_logs(market, symbol, created_at DESC);
