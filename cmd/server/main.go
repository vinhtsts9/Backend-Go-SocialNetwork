package main

import (
	"fmt"
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

	// Start Gin server
	go func() {
		if err := r.Run(":8080"); err != nil {
			log.Fatalf("Failed to start Gin server: %v", err)
		}
	}()
	manager := websocket.NewConnectionManager()

	// Khởi chạy ConnectionManager trong một goroutine
	go manager.Run()

	// Cấu hình route để xử lý kết nối WebSocket
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.HandleConnections(manager, w, r)
	})

	// Chạy server HTTP
	port := "8888"
	fmt.Printf("Server Websocket đang chạy tại http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
