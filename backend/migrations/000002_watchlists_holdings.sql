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
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_watchlists_user_id ON watchlists(user_id);
CREATE INDEX IF NOT EXISTS idx_holdings_user_id ON holdings(user_id);
