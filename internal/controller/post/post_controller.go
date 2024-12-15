package post

import (
	"go-ecommerce-backend-api/m/v2/global"
	model "go-ecommerce-backend-api/m/v2/internal/models"
	"go-ecommerce-backend-api/m/v2/internal/service"
	"go-ecommerce-backend-api/m/v2/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// management post
var Post = new(cPost)

type cPost struct{}

// CreatePost
// @Summary      Create a new post
// @Description  Create a new post for the user
// @Tags         post management
// @Accept       json
// @Produce      json
// @Param        payload body model.CreatePostInput true "Post Payload"
// @Success      201  {object}  response.ResponseData
// @Failure      400  {object}  response.ErrorResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /post/create [post]
func (c *cPost) CreatePost(ctx *gin.Context) {
	var params model.CreatePostInput
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}
	codeRs, dataRs, err := service.NewPost().CreatePost(ctx, &params)
	if err != nil {
		global.Logger.Error("Error creating post", zap.Error(err))
		response.ErrorResponse(ctx, response.ErrCodeInternal, err.Error())
		return
	}
	response.SuccessResponse(ctx, codeRs, dataRs)
}

// UpdatePost
// @Summary      Update a post
// @Description  Update a post by its ID
// @Tags         post management
// @Accept       json
// @Produce      json
// @Param        id   path     int  true  "Post ID"
// @Param        payload body model.UpdatePostInput true "Updated Post Data"
// @Success      200  {object}  response.ResponseData
// @Failure      400  {object}  response.ErrorResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /post/{id} [patch]
func (c *cPost) UpdatePost(ctx *gin.Context) {
	id := ctx.Param("id")
	var params model.UpdatePostInput
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}
	codeRs, dataRs, err := service.NewPost().UpdatePost(ctx, id, &params)
	if err != nil {
		global.Logger.Error("Error updating post", zap.Error(err))
		response.ErrorResponse(ctx, response.ErrCodeInternal, err.Error())
		return
	}
	response.SuccessResponse(ctx, codeRs, dataRs)
}

// DeletePost
// @Summary      Delete a post
// @Description  Delete a post by its ID
// @Tags         post management
// @Accept       json
// @Produce      json
// @Param        id   path     int  true  "Post ID"
// @Success      204  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /post/{id} [delete]
func (c *cPost) DeletePost(ctx *gin.Context) {
	id := ctx.Param("id")
	codeRs, err := service.NewPost().DeletePost(ctx, id)
	if err != nil {
		global.Logger.Error("Error deleting post", zap.Error(err))
		response.ErrorResponse(ctx, response.ErrCodeInternal, err.Error())
		return
	}
	response.SuccessResponse(ctx, codeRs, nil)
}
