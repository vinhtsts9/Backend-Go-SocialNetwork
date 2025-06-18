package user

import (
	"go-ecommerce-backend-api/m/v2/internal/controller/post"
	middleware "go-ecommerce-backend-api/m/v2/internal/middlewares"

	"github.com/gin-gonic/gin"
)

type PostRouter struct {
}

func (p *PostRouter) InitPostRouter(Router *gin.RouterGroup) {
	postRouterPrivate := Router.Group("/post")
	postRouterPrivate.Use(middleware.AuthenMiddleware())
	{
		postRouterPrivate.POST("/create", post.Post.CreatePost)
		postRouterPrivate.PATCH("/:id", post.Post.UpdatePost)
		postRouterPrivate.DELETE("/:id", post.Post.DeletePost)
	}
}
