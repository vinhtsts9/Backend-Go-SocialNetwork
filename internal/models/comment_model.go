package model

type CommentInput struct {
	PostId          uint64 `json:"post_id"`
	UserId          uint64 `json:"user_id"`
	CommentContent  string `json:"comment_content"`
	CommentParentId int32  `json:"comment_parentId"`
}

type ListCommentInput struct {
	PostId          uint64 `json:"post_id"`
	CommentParentId int32  `json:"comment_parentId"`
}

type ListCommentOutput struct {
	Id              int32
	PostId          uint64
	UserId          uint64
	CommentContent  string
	CommentLeft     int32
	CommentRight    int32
	CommentParentId int32
	Isdeleted       bool
}

type DeleteCommentInput struct {
	Id     int32
	PostId uint64
}
