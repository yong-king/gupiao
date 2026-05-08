ALTER TABLE stock_assistant_messages
    ADD COLUMN IF NOT EXISTS session_id TEXT NOT NULL DEFAULT 'default';

CREATE INDEX IF NOT EXISTS idx_stock_assistant_messages_session_symbol
    ON stock_assistant_messages(user_id, session_id, market, symbol, created_at DESC);
