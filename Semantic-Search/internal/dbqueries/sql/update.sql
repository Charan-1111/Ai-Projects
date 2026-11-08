-- name: UpdateDocument
UPDATE documents
SET
    title = $1,
    content = $2,
    category = $3,
    source = $4,
    metadata = $5,
    embedding = $6,
    updated_at = NOW()
WHERE id = $7;

-- name: UpdateIndexingStatus
UPDATE documents
SET indexing_status = 'true'
WHERE id = $1;
