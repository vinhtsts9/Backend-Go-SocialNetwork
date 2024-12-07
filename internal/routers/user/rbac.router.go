package user

import (
	"go-ecommerce-backend-api/m/v2/internal/middlewares"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RbacRouter struct {
}

func (c *RbacRouter) InitRbacRouter(Router *gin.RouterGroup) {
	userRouterPublic := Router.Group("/rbac")
	{
		userRouterPublic.GET("/data1", func(c *gin.Context) {
			c.Set("user", "alice") // Gán user vào context
			userRouterPublic.Use(middlewares.CasbinMiddleware())
			c.JSON(http.StatusOK, gin.H{"message": "Data 1 accessed"})
		})

		userRouterPublic.GET("/data2", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Data 2 accessed"})
		})
	}
}
