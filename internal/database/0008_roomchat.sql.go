// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: 0008_roomchat.sql

package database

import (
	"context"
	"database/sql"
)

const addMemberToRoomChat = `-- name: AddMemberToRoomChat :exec
insert into room_members (user_id, room_id)
values (?,?)
`

type AddMemberToRoomChatParams struct {
	UserID uint64
	RoomID int64
}

func (q *Queries) AddMemberToRoomChat(ctx context.Context, arg AddMemberToRoomChatParams) error {
	_, err := q.db.ExecContext(ctx, addMemberToRoomChat, arg.UserID, arg.RoomID)
	return err
}

const createRoomChat = `-- name: CreateRoomChat :exec
insert into chat_rooms(name, is_group, admin_id, avatar_url)
values (?,?,?,?)
`

type CreateRoomChatParams struct {
	Name      string
	IsGroup   bool
	AdminID   uint64
	AvatarUrl string
}

func (q *Queries) CreateRoomChat(ctx context.Context, arg CreateRoomChatParams) error {
	_, err := q.db.ExecContext(ctx, createRoomChat,
		arg.Name,
		arg.IsGroup,
		arg.AdminID,
		arg.AvatarUrl,
	)
	return err
}

const deleteMemberFromRoomChat = `-- name: DeleteMemberFromRoomChat :exec
delete from room_members
where user_id = ? and room_id = ?
`

type DeleteMemberFromRoomChatParams struct {
	UserID uint64
	RoomID int64
}

func (q *Queries) DeleteMemberFromRoomChat(ctx context.Context, arg DeleteMemberFromRoomChatParams) error {
	_, err := q.db.ExecContext(ctx, deleteMemberFromRoomChat, arg.UserID, arg.RoomID)
	return err
}

const getMemberGroup = `-- name: GetMemberGroup :many
select ui.user_nickname, ui.user_avatar
from user_info ui
join room_members rb on ui.user_id = rb.user_id
where rb.room_id = ?
`

type GetMemberGroupRow struct {
	UserNickname sql.NullString
	UserAvatar   sql.NullString
}

func (q *Queries) GetMemberGroup(ctx context.Context, roomID int64) ([]GetMemberGroupRow, error) {
	rows, err := q.db.QueryContext(ctx, getMemberGroup, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetMemberGroupRow
	for rows.Next() {
		var i GetMemberGroupRow
		if err := rows.Scan(&i.UserNickname, &i.UserAvatar); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRoomByUserId = `-- name: GetRoomByUserId :many
SELECT cr.id, cr.name, cr.is_group, cr.admin_id, cr.avatar_url, cr.created_at 
FROM chat_rooms cr
JOIN room_members rb ON cr.id = rb.room_id
WHERE rb.user_id = ?
`

func (q *Queries) GetRoomByUserId(ctx context.Context, userID uint64) ([]ChatRoom, error) {
	rows, err := q.db.QueryContext(ctx, getRoomByUserId, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ChatRoom
	for rows.Next() {
		var i ChatRoom
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.IsGroup,
			&i.AdminID,
			&i.AvatarUrl,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
