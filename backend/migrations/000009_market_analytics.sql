ALTER TABLE price_snapshots ADD COLUMN IF NOT EXISTS name TEXT NOT NULL DEFAULT '';
ALTER TABLE price_snapshots ADD COLUMN IF NOT EXISTS open NUMERIC NOT NULL DEFAULT 0;
ALTER TABLE price_snapshots ADD COLUMN IF NOT EXISTS high NUMERIC NOT NULL DEFAULT 0;
ALTER TABLE price_snapshots ADD COLUMN IF NOT EXISTS low NUMERIC NOT NULL DEFAULT 0;

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

CREATE INDEX IF NOT EXISTS idx_daily_price_changes_symbol_date ON daily_price_changes(market, symbol, date DESC);
