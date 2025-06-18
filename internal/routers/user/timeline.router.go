// package user

// import (
// 	"go-ecommerce-backend-api/m/v2/internal/controller/comment"

// 	"github.com/gin-gonic/gin"
// )

// type CommentRouter struct {
// }

// func (cr *CommentRouter) InitCommentRouter(Router *gin.RouterGroup) {
// 	commentRouterPrivate := Router.Group("/timeline")

// 	commentRouterPrivate.POST("/create", comment.Comment.CreateComment)
// 	commentRouterPrivate.GET("/:post_id/:comment_parentId", comment.Comment.ListComment) // Nhận tham số động
// 	commentRouterPrivate.DELETE("/:id/:post_id", comment.Comment.DeleteComment)
// }

package user

import (
	"go-ecommerce-backend-api/m/v2/internal/controller/timeline"
	"go-ecommerce-backend-api/m/v2/internal/middlewares"

	"github.com/gin-gonic/gin"
)

type TimelineRouter struct {
}

func (tr *TimelineRouter) InitTimelineRouter(Router *gin.RouterGroup) {
	timelineRouterPrivate := Router.Group("/timeline")
	timelineRouterPrivate.Use(middlewares.AuthenMiddleware())
	{
		timelineRouterPrivate.GET("/:post_id", timeline.Timeline.GetPost)
		timelineRouterPrivate.GET("/all", timeline.Timeline.GetAllPost)
	}

}
