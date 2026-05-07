CREATE TABLE IF NOT EXISTS analysis_reports (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id),
    report_date DATE NOT NULL,
    summary TEXT NOT NULL,
    risk_points JSONB NOT NULL DEFAULT '[]',
    needs_confirmation JSONB NOT NULL DEFAULT '[]',
    data_time TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_analysis_reports_user_date ON analysis_reports(user_id, report_date DESC);
