-- VECTOR EXTENSION
CREATE EXTENSION IF NOT EXISTS vector;

-- CHUNKS TABLE
CREATE TABLE IF NOT EXISTS chunks (
    id TEXT PRIMARY KEY,
    owner TEXT NOT NULL,
    repo TEXT NOT NULL,
    filename TEXT NOT NULL,
    content TEXT NOT NULL,
    embedding vector(384)
);

CREATE INDEX IF NOT EXISTS idx_chunks_embedding ON chunks USING hnsw (embedding vector_cosine_ops);

-- HISTORY TABLE
CREATE TABLE IF NOT EXISTS chat_history (
    id SERIAL PRIMARY KEY,
    session_id VARCHAR(255),
    role VARCHAR(50),
    content TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);