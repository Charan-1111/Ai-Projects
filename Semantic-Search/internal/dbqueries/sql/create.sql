-- name: CreateVectorExtension
CREATE EXTENSION IF NOT EXISTS vector;

-- name: CreateDocumentsTable
CREATE TABLE IF NOT EXISTS documents (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    category TEXT,
    source TEXT,
    metadata JSONB,
    embedding VECTOR(1536) NOT NULL,
    indexing_status TEXT NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- name: CreateChunksTable
CREATE TABLE IF NOT EXISTS chunks (
    id UUID PRIMARY KEY,
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    chunk_index INTEGER NOT NULL,
    start_position INTEGER,
    end_position INTEGER,
    embedding VECTOR(1536) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
