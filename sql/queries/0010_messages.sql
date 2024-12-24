
-- name: GetChatHistory :many
select user_nickname, message_context, message_type, is_pinned, created_at from messages
where room_id = ? order by created_at asc ;

-- name: SetChatHistory :exec
insert into messages(user_nickname, message_context, message_type,room_id) values (?,?,?,?) 