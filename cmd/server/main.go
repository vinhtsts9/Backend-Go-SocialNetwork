package main

import (
	"go-ecommerce-backend-api/m/v2/internal/initialize"
	websocket "go-ecommerce-backend-api/m/v2/third_party/ws"
	"log"
	"net/http"

	_ "go-ecommerce-backend-api/m/v2/cmd/swag/docs"

	"github.com/gin-gonic/gin"
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

// @host      localhost:8080
// @BasePath  /v1/2024
// @schema http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// Gin for API
	r := initialize.Run()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/checkStatus", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	cm := websocket.NewConnectionManager()
	cm.Run()

	// Route WebSocket được xử lý trong Gin server

	r.GET("/ws", func(c *gin.Context) {
		websocket.HandleConnections(c.Writer, c.Request, cm)
	})

	log.Println("Starting server on :8080")

	go func() {
		log.Println("Starting server on :8080")
		if err := r.Run(":8080"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	select {}
}
