CREATE TABLE IF NOT EXISTS price_snapshots (
    id TEXT PRIMARY KEY,
    market TEXT NOT NULL,
    symbol TEXT NOT NULL,
    price NUMERIC NOT NULL,
    previous_close NUMERIC NOT NULL,
    change_percent NUMERIC NOT NULL,
    volume BIGINT NOT NULL,
    source TEXT NOT NULL,
    data_time TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
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

CREATE INDEX IF NOT EXISTS idx_price_snapshots_symbol_time ON price_snapshots(market, symbol, data_time DESC);
CREATE INDEX IF NOT EXISTS idx_refresh_jobs_user_id ON refresh_jobs(user_id);
