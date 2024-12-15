package service

import (
	model "go-ecommerce-backend-api/m/v2/internal/models"

	"github.com/gin-gonic/gin"
)

type (
	TimelineInterface interface {
		GetAllPosts(ctx *gin.Context, userId int64) (codeRs int, data []model.Post, err error)
		GetPostById(ctx *gin.Context, postId string) (codeRs int, data model.Post, err error)
	}
)

var localTimelineInterface TimelineInterface

func InitTimelineInterface(i TimelineInterface) {
	localTimelineInterface = i
}

func NewTimelineInterface() TimelineInterface {
	if localTimelineInterface == nil {
		panic("Failed to init timeline interface")
	}
	return localTimelineInterface
}
