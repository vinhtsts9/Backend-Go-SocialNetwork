package model

type CreateRoom struct {
	NameRoom  string `json:"name"`
	IsGroup   bool   `json:"is_group"`
	AdminId   uint64 `json:"admin_id"`
	AvatarUrl string `json:"avatar_url"`
}
