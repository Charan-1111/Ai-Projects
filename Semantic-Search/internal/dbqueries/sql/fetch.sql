-- name: FetchSimilarDocuments
SELECT
    id,
    title,
    content,
    category,
    1 - (embedding <=> $1::vector) AS similarity
FROM documents
ORDER BY embedding <=> $1::vector
LIMIT $2;

-- name: FetchChunkedDocuments
SELECT
    c.id,
    c.document_id,
    c.content,
    c.chunk_index,
    1 - (c.embedding <=> $1::vector) AS similarity,
    d.title,
    d.category,
    d.metadata
FROM chunks c
JOIN documents d ON d.id = c.document_id
WHERE 1 - (c.embedding <=> $1::vector) >= $3
ORDER BY c.embedding <=> $1
LIMIT $2;

-- name: FetchFilteredChunkedDocuments
SELECT
    c.id,
    c.document_id,
    c.content,
    c.chunk_index,
    1 - (c.embedding <=> $1::vector) AS similarity,
    d.title,
    d.category,
    d.metadata
FROM chunks c
JOIN documents d ON d.id = c.document_id
WHERE ($2 = '' OR d.category = $2)
    AND ($3 = '' OR d.metadata->>'difficulty' = $3)
    AND 1 - (c.embedding <=> $1::vector) >= $5
ORDER BY c.embedding <=> $1
LIMIT $4;

-- name: FetchDocumentByID
SELECT
    id,
    title,
    content,
    category,
    source,
    metadata
FROM documents
WHERE id = $1;
