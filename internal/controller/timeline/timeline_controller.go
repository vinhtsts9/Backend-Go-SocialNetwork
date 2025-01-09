package timeline

import (
	"go-ecommerce-backend-api/m/v2/global"
	"go-ecommerce-backend-api/m/v2/internal/service"
	"go-ecommerce-backend-api/m/v2/package/utils/auth"
	"go-ecommerce-backend-api/m/v2/response"

	"github.com/gin-gonic/gin"
)

type cTimeline struct {
}

var Timeline = new(cTimeline)

func (c *cTimeline) GetAllPost(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	userInfo := auth.GetUserInfoFromToken(token)
	global.Logger.Sugar().Info(userInfo.UserID)
	codeRs, data, err := service.NewTimelineInterface().GetAllPosts(ctx, int64(userInfo.UserID))
	if err != nil {
		global.Logger.Sugar().Error("get post error ", err)
		response.ErrorResponse(ctx, response.ErrCodePostFailed, "get posts failed")
	}
	global.Logger.Sugar().Info("Time of post", data)
	response.SuccessResponse(ctx, codeRs, data)
}

func (c *cTimeline) GetPost(ctx *gin.Context) {
	postId := ctx.Param("post_id")
	codeRs, data, err := service.NewTimelineInterface().GetPostById(ctx, postId)
	if err != nil {
		global.Logger.Sugar().Error("get post error ", err)
		response.ErrorResponse(ctx, response.ErrCodePostFailed, "get post failed")
	}
	response.SuccessResponse(ctx, codeRs, data)
}
