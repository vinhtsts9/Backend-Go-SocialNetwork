package user

import (
	"go-ecommerce-backend-api/m/v2/internal/controller/comment"

	"github.com/gin-gonic/gin"
)

type CommentRouter struct {
}

func (cr *CommentRouter) InitCommentRouter(Router *gin.RouterGroup) {
	commentRouterPrivate := Router.Group("/comment")

	commentRouterPrivate.POST("/create", comment.Comment.CreateComment)
	commentRouterPrivate.GET("/:post_id/:comment_parentId", comment.Comment.ListComment) // Nhận tham số động
	commentRouterPrivate.DELETE("/:id/:post_id", comment.Comment.DeleteComment)
}
