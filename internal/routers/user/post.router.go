package user

// import (
// 	"go-ecommerce-backend-api/m/v2/internal/middlewares"

// 	"github.com/gin-gonic/gin"
// )

// type PostRouter struct {
// }

// func (p *PostRouter) InitPostRouter(Router *gin.RouterGroup) {
// 	postRouterPublic := Router.Group("/post")
// 	{
// 		postRouterPublic.GET("/list", post.Post.GetPosts)
// 		postRouterPublic.Get(":/id", post.Post.GetPostById)
// 	}

// 	postRouterPrivate := Router.Group("/post")
// 	postRouterPrivate.Use(middlewares.AuthenMiddleware())
// 	{
// 		postRouterPrivate.Post("/create", post.Post.CreatePost)
// 		postRouterPrivate.PATCH("/:id", post.Post.UpdatePost)
// 		postRouterPrivate.DELETE("/:id", post.Post.DeletePost)
// 	}
// }
