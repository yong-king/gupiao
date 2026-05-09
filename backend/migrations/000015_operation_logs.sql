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

CREATE INDEX IF NOT EXISTS idx_operation_logs_user_created
    ON operation_logs(user_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_operation_logs_symbol_created
    ON operation_logs(market, symbol, created_at DESC);
