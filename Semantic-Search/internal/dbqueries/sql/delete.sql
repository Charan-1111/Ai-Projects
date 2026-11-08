-- name: DeleteDocument
DELETE FROM documents
WHERE id = $1;

-- name: DeleteChunks
DELETE FROM chunks
WHERE document_id = $1;
