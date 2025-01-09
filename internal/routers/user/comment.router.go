package user

import (
	"go-ecommerce-backend-api/m/v2/internal/controller/comment"
	"go-ecommerce-backend-api/m/v2/internal/middlewares"

	"github.com/gin-gonic/gin"
)

type CommentRouter struct {
}

func (cr *CommentRouter) InitCommentRouter(Router *gin.RouterGroup) {
	commentRouterPrivate := Router.Group("/comment")
	commentRouterPrivate.Use(middlewares.AuthenMiddleware())
	{

		commentRouterPrivate.POST("/create", comment.Comment.CreateComment)
		commentRouterPrivate.GET("/:post_id/:comment_parentId", comment.Comment.ListComment)

		commentRouterPrivate.GET("/:post_id/root", comment.Comment.ListCommentRoot) // Nhận tham số động
		commentRouterPrivate.DELETE("/:id/:post_id", comment.Comment.DeleteComment)
	}
}
