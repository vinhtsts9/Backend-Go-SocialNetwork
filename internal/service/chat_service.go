package service

import (
	model "go-ecommerce-backend-api/m/v2/internal/models"

	"github.com/gin-gonic/gin"
)

type (
	IChat interface {
		CreateRoom(ctx *gin.Context, RoomModel *model.CreateRoom) (codeRs int, err error)
		GetChatHistory(ctx *gin.Context, roomId int) (codeRs int, err error, rs []model.ModelChat)
		SetChatHistory(ctx *gin.Context, model *model.ModelChat)
		GetRoomChatByUserId(ctx *gin.Context, userId uint64) (codeRs int, rs []model.CreateRoom, err error)
	}
)

var localIChat IChat

func NewIChat() IChat {
	if localIChat == nil {
		panic("Init IChat failed ")
	}
	return localIChat
}

func InitIChat(i IChat) {
	localIChat = i
}
