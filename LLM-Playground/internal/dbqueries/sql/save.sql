-- name: SaveConversation
INSERT INTO conversations (id, user_id)
VALUES ($1, $1)
ON CONFLICT (id) DO NOTHING;

-- name: SaveMessage
INSERT INTO messages (id, conversation_id, role, content)
VALUES ($1, $2, $3, $4);
