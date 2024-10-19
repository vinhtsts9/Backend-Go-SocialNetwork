package initialize

import (
	"go-ecommerce-backend-api/m/v2/global"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Run() *gin.Engine {
	LoadConfig()
	InitLogger()
	global.Logger.Info("Config ok", zap.String("ok", "success"))
	InitMysqlC()
	InitServiceInterface()
	InitRedis()
	InitKafKa()
	r := InitRouter()
	return r
}
