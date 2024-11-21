package initialize

import (
	"go-ecommerce-backend-api/m/v2/global"
	"go-ecommerce-backend-api/m/v2/package/logger"
)

func InitLogger() {
	lz := logger.NewLogger(global.Config.Logger)
	global.Logger = lz.GetZapLogger()
}
