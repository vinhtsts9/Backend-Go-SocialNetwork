package comment

import (
	"go-ecommerce-backend-api/m/v2/global"
	model "go-ecommerce-backend-api/m/v2/internal/models"
	"go-ecommerce-backend-api/m/v2/internal/service"
	"go-ecommerce-backend-api/m/v2/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type cComment struct{}

var Comment = new(cComment)

// Createcomment
// @Summary      Create a new comment
// @Description  Create a new comment for the user
// @Tags         comment management
// @Accept       json
// @Produce      json
// @Param        payload body model.CreateCommentInput true "comment Payload"
// @Success      201  {object}  response.ResponseData
// @Failure      400  {object}  response.ErrorResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /comment/create [post]

func (c *cComment) CreateComment(ctx *gin.Context) {
	var params model.CreateCommentInput
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeComment, err.Error())
		return
	}

	codeRs, err := service.NewICommnet().CreateComment(ctx, &params)
	if err != nil {
		global.Logger.Sugar().Error("Create comment error:", err)
		response.ErrorResponse(ctx, response.ErrCodeComment, err.Error())
		return
	}
	response.SuccessResponse(ctx, codeRs, 1)
}

// Listcomment
// @Summary      List comments
// @Description  Get a list of comments for a specific post and parent comment
// @Tags         comment management
// @Accept       json
// @Produce      json
// @Param        post_id          path      int  true  "Post ID"
// @Param        comment_parentId path      int  true  "Parent Comment ID"
// @Success      200  {object}  response.ResponseData
// @Failure      400  {object}  response.ErrorResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /comment/{post_id}/{comment_parentId} [get]

func (c *cComment) ListComment(ctx *gin.Context) {
	postID := ctx.Param("post_id")
	commentParentId := ctx.Param("comment_parentId")

	postId, err := strconv.ParseUint(postID, 10, 64)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeComment, "Invalid post_id")
	}
	commentParentID, err := strconv.ParseInt(commentParentId, 10, 32)
	if err != nil && ctx.Param("comment_parentId") != "" {
		response.ErrorResponse(ctx, response.ErrCodeComment, "Invalid comment_parentId")
	}
	params := model.ListCommentInput{
		PostId:          postId,
		CommentParentId: int32(commentParentID),
	}
	codeRs, err, data := service.NewICommnet().ListComments(&params)
	if err != nil {
		global.Logger.Sugar().Error("List comment error ", err)
		response.ErrorResponse(ctx, response.ErrCodeComment, err.Error())
		return
	}
	response.SuccessResponse(ctx, codeRs, data)
}

// Deletecomment
// @Summary      Delete comment
// @Description  Delete comment for the user
// @Tags         comment management
// @Accept       json
// @Produce      json
// @Param        id  path int true "Id"
// @Param        post_id  path int true "PostId"
// @Success      201  {object}  response.ResponseData
// @Failure      400  {object}  response.ErrorResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /comment/{id}/{post_id} [delete]
func (c *cComment) DeleteComment(ctx *gin.Context) {
	ID := ctx.Param("id")
	PostID := ctx.Param("post_id")

	Id, err := strconv.ParseInt(ID, 10, 32)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeComment, "Invalid post_id")
	}
	PostId, err := strconv.ParseUint(PostID, 10, 64)
	if err != nil && ctx.Param("post_id") != "" {
		response.ErrorResponse(ctx, response.ErrCodeComment, "Invalid comment_parentId")
	}
	params := model.DeleteCommentInput{
		Id:     int32(Id),
		PostId: PostId,
	}
	codeRs, err, data := service.NewICommnet().DeleteComment(&params)
	if err != nil {
		global.Logger.Sugar().Error("Delete Comment Failed", err)
		response.ErrorResponse(ctx, response.ErrCodeComment, err.Error())
		return
	}
	response.SuccessResponse(ctx, codeRs, data)
}
