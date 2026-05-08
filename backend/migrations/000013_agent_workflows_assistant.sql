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
    market TEXT NOT NULL,
    symbol TEXT NOT NULL,
    question TEXT NOT NULL,
    answer TEXT NOT NULL,
    context_summary TEXT NOT NULL DEFAULT '',
    rag_document_ids JSONB NOT NULL DEFAULT '[]'::jsonb,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_agent_workflow_jobs_user_created ON agent_workflow_jobs(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_agent_workflow_steps_job ON agent_workflow_steps(job_id, started_at);
CREATE INDEX IF NOT EXISTS idx_stock_assistant_messages_user_symbol ON stock_assistant_messages(user_id, market, symbol, created_at DESC);
