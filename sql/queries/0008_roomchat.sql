-- name: AddMemberToRoomChat :exec
insert into room_members (user_id, room_id)
values (?,?);

-- name: DeleteMemberFromRoomChat :exec 
delete from room_members
where user_id = ? and room_id = ?;

-- name: CreateRoomChat :exec
insert into room_chats(name, is_group, admin_id, avatar_url)
values (?,?,?,?);

-- name: GetRoomByUserId :many
SELECT rc.* 
FROM room_chats rc
JOIN room_members rb ON rc.id = rb.room_id
WHERE rb.user_id = ?;
