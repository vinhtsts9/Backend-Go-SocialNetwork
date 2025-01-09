package user

import (
	"go-ecommerce-backend-api/m/v2/internal/controller/chat"
	"go-ecommerce-backend-api/m/v2/internal/middlewares"

	"github.com/gin-gonic/gin"
)

type ChatRouter struct {
}

func (c *ChatRouter) InitChatRouter(Router *gin.RouterGroup) {
	chatRouterGroup := Router.Group("/chat")
	chatRouterGroup.Use(middlewares.AuthenMiddleware())
	{
		chatRouterGroup.POST("/createRoom", chat.Chat.CreateRoom)
		chatRouterGroup.GET("/get_history/:room_id", chat.Chat.GetChatHistory)
		chatRouterGroup.GET("/get_user_nickname", chat.Chat.GetUserNickName)
		chatRouterGroup.GET("/get_room_by_userId", chat.Chat.GetRoomByUserId)
	}
	chatRouterGroup.Use(middlewares.CasbinMiddleware())
	{
		chatRouterGroup.GET("/get_member/:room_id", chat.Chat.GetMemberGroup)
		chatRouterGroup.DELETE("/delete_member/:room_id/:user_id", chat.Chat.DeleteMemberFromGroup)
	}
}
