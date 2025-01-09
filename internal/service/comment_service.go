package service

import (
	model "go-ecommerce-backend-api/m/v2/internal/models"

	"github.com/gin-gonic/gin"
)

type (
	IComment interface {
		CreateComment(ctx *gin.Context, model *model.CreateCommentInput, userId uint64) (codeRs int, rs model.ListCommentOutput, err error)
		ListComments(*model.ListCommentInput) (codeRs int, data []model.ListCommentOutput, err error)
		ListCommentRoot(ctx *gin.Context, postId uint64) (codeRs int, data []model.ListCommentOutput, err error)
		DeleteComment(*model.DeleteCommentInput) (codeRs int, err error, Rs bool)
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
