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

CREATE INDEX IF NOT EXISTS idx_rag_documents_symbol_created ON rag_documents(market, symbol, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_llm_analysis_runs_user_symbol ON llm_analysis_runs(user_id, market, symbol, created_at DESC);
