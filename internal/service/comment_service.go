package service

import (
	"go-ecommerce-backend-api/m/v2/internal/database"
	model "go-ecommerce-backend-api/m/v2/internal/models"
)

type (
	IComment interface {
		CreateComment(*model.CommentInput) (codeRs int, err error)
		ListComments(*model.ListCommentInput) (codeRs int, err error, data []database.Comment)
		DeleteComment(*model.DeleteCommentInput) (bool, error)
	}
)

var localICommnet IComment

func InitIComment(i IComment) {
	localICommnet = i
}

func NewICommnet() IComment {
	if localICommnet == nil {
		panic("Init ICommnet failed")
	}
	return localICommnet
}
