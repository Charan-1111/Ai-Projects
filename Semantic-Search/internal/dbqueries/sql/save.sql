-- name: SaveDocument
INSERT INTO documents (
    id,
    title,
    content,
    category,
    source,
    metadata,
    embedding
)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: SaveChunk
INSERT INTO chunks (
    id,
    document_id,
    content,
    chunk_index,
    start_position,
    end_position,
    embedding
)
VALUES ($1, $2, $3, $4, $5, $6, $7);
