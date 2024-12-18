package initialize

import (
	"go-ecommerce-backend-api/m/v2/global"
	"go-ecommerce-backend-api/m/v2/internal/middlewares"
	"go-ecommerce-backend-api/m/v2/internal/routers"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	var r *gin.Engine
	if global.Config.Server.Mode == "dev" {
		gin.SetMode(gin.DebugMode)
		gin.ForceConsoleColor()
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()

	}
	// middlewares
	// r.Use() logging
	// r.Use() cors
	r.Use(middlewares.CORSMiddleware())
	// r.Use() limiter global
	manageRouter := routers.RouterGroupApp.Manage
	userRouter := routers.RouterGroupApp.User

	MainGroup := r.Group("/v1/2024")
	{
		// MainGroup.GET("/checkStatus") //checking monitor

	}
	{
		userRouter.InitUserRouter(MainGroup)
		userRouter.InitProductRouter(MainGroup)
		userRouter.InitPostRouter(MainGroup)
		userRouter.InitRbacRouter(MainGroup)
		userRouter.InitChatRouter(MainGroup)
		userRouter.InitCommentRouter(MainGroup)
		userRouter.InitTimelineRouter(MainGroup)
	}
	{
		manageRouter.InitAdminRouter(MainGroup)
		manageRouter.InitUserRouter(MainGroup)
	}
	return r
}
