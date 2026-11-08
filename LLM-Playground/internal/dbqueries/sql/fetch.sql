-- name: FetchMessages
SELECT role, content
FROM messages
WHERE conversation_id = $1
ORDER BY created_at ASC;
