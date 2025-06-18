-- name: GetChatHistory :many
SELECT user_nickname, message_context, message_type, is_pinned, created_at 
FROM messages
WHERE room_id = ? 
ORDER BY created_at ASC
LIMIT 10;


-- name: SetChatHistory :exec
insert into messages(user_nickname, message_context, message_type,room_id) values (?,?,?,?) 