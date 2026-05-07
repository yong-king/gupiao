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

CREATE INDEX IF NOT EXISTS idx_alert_rules_user_symbol ON alert_rules(user_id, market, symbol);
CREATE INDEX IF NOT EXISTS idx_alert_events_user_created ON alert_events(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_notifications_user_created ON notifications(user_id, created_at DESC);
