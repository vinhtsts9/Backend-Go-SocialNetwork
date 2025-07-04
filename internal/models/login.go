package model

import (
	"database/sql"
	"time"
)

type UserInfo struct {
	UserID             int        `json:"UserID"`
	UserAccount        string     `json:"UserAccount"`
	UserNickname       NullString `json:"UserNickname"`
	UserAvatar         NullString `json:"UserAvatar"`
	UserState          int        `json:"UserState"`
	UserMobile         NullString `json:"UserMobile"`
	UserGender         NullInt16  `json:"UserGender"`
	UserBirthday       NullTime   `json:"UserBirthday"`
	UserEmail          NullString `json:"UserEmail"`
	UserIsAuthencation int        `json:"UserIsAuthencation"`
	CreatedAt          NullTime   `json:"CreatedAt"`
	UpdatedAt          NullTime   `json:"UpdatedAt"`
}
type UserBase struct {
	UserID         int32
	UserLogoutTime sql.NullTime
	UserState      uint8
}
type UserSearch struct {
	UserNickname string `json:"user_nickname"`
	UserAvatar   string `json:"user_avatar"`
}

// Các kiểu phụ trợ cho Null* (giống như trong database/sql)
type NullString struct {
	String string `json:"String"`
	Valid  bool   `json:"Valid"`
}

type NullInt16 struct {
	Int16 int16 `json:"Int16"`
	Valid bool  `json:"Valid"`
}

type NullTime struct {
	Time  time.Time `json:"Time"`
	Valid bool      `json:"Valid"`
}

type Message struct {
	MessageId      uint32 `json:"id"`
	RoomId         uint32 `json:"room_id"`
	SenderId       uint64 `json:"sender_id"`
	MessageContext string `json:"message_context"`
	MessageType    string `json:"message_type"`
	IsPinned       uint8  `json:"is_pinned"`
	IsAnnouncement uint8  `json:"is_announcement"`
}
type RegisterInput struct {
	VerifyKey     string `json:"verify_key"`
	VerifyType    int    `json:"verify_type"`
	VerifyPurpose string `json:"verify_purpose"`
}

type VerifyInput struct {
	VerifyKey  string `json:"verify_key"`
	VerifyCode string `json:"verify_code"`
}

type VerifyOTPOutput struct {
	Token   string `json:"token"`
	UserId  string `json:"userId"`
	Message string `json:"message"`
}

type UpdatePasswordRegisterInput struct {
	UserToken    string    `json:"user_token"`
	UserPassword string    `json:"user_password"`
	UserNickname string    `json:"user_nickname"`
	UserAvatar   string    `json:"user_avatar"`
	UserMobile   string    `json:"user_mobile"`
	UserGender   int16     `json:"user_gender"`
	UserBirthday time.Time `json:"user_birthday"`
}
type LoginInput struct {
	UserAccount  string `json:"user_account"`
	UserPassword string `json:"user_password"`
}
type LoginOutPut struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
type LogoutInput struct {
	TokenString  string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

// two factor authentication
type SetupTwoFactorAuthInput struct {
	UserId            uint32 `json:"user_id"`
	TwoFactorAuthType string `json:"two_factor_auth_type"`
	TwoFactorEmail    string `json:"two_factor_email"`
}

type TwoFactorVerificationInput struct {
	UserId        uint32 `json:"user_id"`
	TwoFactorCode string `json:"two_factor_code"`
}

// Post
// Model dùng để tạo mới Post
type CreatePostInput struct {
	UserId       uint64      `json:"user_id"`
	UserNickname NullString  `json:"user_nickname"`
	Title        string      `json:"title"`
	ImagePaths   interface{} `json:"image_paths"` // JSON content in string format
	Privacy      string      `json:"privacy_mode"`
	Metadata     string      `json:"metadata"` // JSON metadata in string format
}

// Deadline implements context.Context.
func (c *CreatePostInput) Deadline() (deadline time.Time, ok bool) {
	panic("unimplemented")
}

// Done implements context.Context.
func (c *CreatePostInput) Done() <-chan struct{} {
	panic("unimplemented")
}

// Err implements context.Context.
func (c *CreatePostInput) Err() error {
	panic("unimplemented")
}

// Value implements context.Context.
func (c *CreatePostInput) Value(key any) any {
	panic("unimplemented")
}

type Post struct {
	ID           uint32   `json:"id"`
	UserId       uint64   `json:"user_id"`
	UserNickname string   `json:"user_nickname"`
	Title        string   `json:"title"`
	ImagePaths   []string `json:"image_paths"`
	CreatedAt    string   `json:"created_at"` // Format time string
	UpdatedAt    string   `json:"updated_at"` // Format time string
	IsPublished  bool     `json:"is_published"`
	Metadata     string   `json:"metadata"`
}

// Model dùng để cập nhật Post
type UpdatePostInput struct {
	ID           uint32      `json:"id"`
	UserId       uint64      `json:"user_id"`
	UserNickname string      `json:"user_nickname"`
	Title        string      `json:"title"`
	ImagePaths   interface{} `json:"image_paths"` // JSON content in string format
	IsPublished  bool        `json:"is_published"`
	Metadata     string      `json:"metadata"` // JSON metadata in string format
}

// Model Post trả về cho người dùng

type CasbinPolicy struct {
	PType string
	V0    string
	V1    string
	V2    string
	V3    string
	V4    string
}
type UpdateAvatar struct {
	UserAvatar string `json:"user_avatar"`
}
