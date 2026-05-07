CREATE TABLE IF NOT EXISTS agent_runs (
    id TEXT PRIMARY KEY,
    agent_version TEXT NOT NULL,
    prompt_version TEXT NOT NULL,
    status TEXT NOT NULL,
    input JSONB NOT NULL,
    output JSONB NOT NULL DEFAULT '{}',
    error TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL
);
