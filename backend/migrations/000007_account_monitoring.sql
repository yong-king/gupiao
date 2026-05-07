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

CREATE INDEX IF NOT EXISTS idx_broker_accounts_user_id ON broker_accounts(user_id);
