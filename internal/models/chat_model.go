package model

import (
	"database/sql"
)

type CreateRoom struct {
	Id        int32  `json:"room_id"`
	NameRoom  string `json:"name"`
	IsGroup   bool   `json:"is_group"`
	AdminId   uint64 `json:"admin_id"`
	AvatarUrl string `json:"avatar_url"`
}
type MessagesMessageType string

const (
	MessagesMessageTypeText  MessagesMessageType = "text"
	MessagesMessageTypeImage MessagesMessageType = "image"
	MessagesMessageTypeVideo MessagesMessageType = "video"
	MessagesMessageTypeFile  MessagesMessageType = "file"
)

// NullMessagesMessageType là kiểu dữ liệu cho message_type

// ModelChat là cấu trúc để giải mã dữ liệu lịch sử chat
type ModelChat struct {
	UserNickname   string              `json:"user_nickname"`
	MessageContext sql.NullString      `json:"message_context"`
	MessageType    MessagesMessageType `json:"message_type"`
	IsPinned       sql.NullBool        `json:"is_pinned"`
	RoomId         sql.NullInt32       `json:"room_id"`
	CreatedAt      sql.NullTime        `json:"created_at"`
}
