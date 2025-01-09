package model

type CreateCommentInput struct {
	PostId          uint64 `json:"post_id"`
	CommentContent  string `json:"comment_content"`
	CommentParentId int32  `json:"comment_parentId"`
	UserNickname    string `json:"user_nickname"`
}

type ListCommentInput struct {
	PostId          uint64 `json:"post_id"`
	CommentParentId int32  `json:"comment_parentId"`
}

type ListCommentOutput struct {
	Id              int32  `json:"comment_id"`
	PostId          uint64 `json:"post_id"`
	UserNickname    string `json:"user_nickname"`
	CommentContent  string `json:"comment_content"`
	CommentLeft     int32  `json:"comment_left"`
	CommentRight    int32  `json:"comment_right"`
	CommentParentId int32  `json:"comment_parentId"`
	ReplyCount      int32  `json:"reply_count"`
	Isdeleted       bool   `json:"isDeleted"`
	CreatedAt       string `json:"created_at"`
}

type DeleteCommentInput struct {
	Id     int32
	PostId uint64
}
