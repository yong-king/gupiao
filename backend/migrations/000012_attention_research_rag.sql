ALTER TABLE holdings ADD COLUMN IF NOT EXISTS attention_level TEXT NOT NULL DEFAULT 'medium';

CREATE TABLE IF NOT EXISTS rag_vectors (
    id TEXT PRIMARY KEY,
    rag_document_id TEXT NOT NULL REFERENCES rag_documents(id),
    provider TEXT NOT NULL,
    model TEXT NOT NULL,
    embedding JSONB NOT NULL DEFAULT '[]'::jsonb,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_rag_vectors_document ON rag_vectors(rag_document_id);
