package initialize

import (
	"go-ecommerce-backend-api/m/v2/global"
	"go-ecommerce-backend-api/m/v2/internal/database"
	"go-ecommerce-backend-api/m/v2/internal/service"
	"go-ecommerce-backend-api/m/v2/internal/service/impl"
)

func InitServiceInterface() {
	queries := database.New(global.Mdbc)
	service.InitUserLogin(impl.NemUserLoginImpl(queries))
	service.InitPost(impl.NewPostImpl(queries))
	service.InitRBACService(impl.NewRbacImpl(queries))
	service.InitPostProcessor(impl.NewPostProcessorImpl(queries))
	service.InitIChat(impl.NewsChat(queries))
}
