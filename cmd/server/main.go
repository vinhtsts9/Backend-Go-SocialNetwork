package main

import (
	"go-ecommerce-backend-api/m/v2/internal/initialize"

	_ "go-ecommerce-backend-api/m/v2/cmd/swag/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger" // swagger embed files
)

// @title           API Documentation Ecommerce Backend Shopdevgo
// @version         1.0.0
// @description     This is a sample server celler server.
// @termsOfService  github.com/Vinhts/GO-MAIN

// @contact.name   Team go
// @contact.url	   github.com/Vinhts/GO-MAIN
// @contact.email  vinhtiensinh17@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:1006
// @BasePath  /v1/2024
// @schema http
func main() {
	r := initialize.Run()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":1006")
}
