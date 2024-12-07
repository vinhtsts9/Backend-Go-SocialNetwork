package chat

import (
	model "go-ecommerce-backend-api/m/v2/internal/models"
	"go-ecommerce-backend-api/m/v2/internal/service"
	"go-ecommerce-backend-api/m/v2/response"

	"github.com/gin-gonic/gin"
)

var Chat = new(cChat)

type cChat struct {
}

// CreateNewRoom
// @Summary      Create a new room
// @Description  Create a new room for the user
// @Tags         chat management
// @Accept       json
// @Produce      json
// @Param        payload body model.CreateRoom true "CreateRoom Payload"
// @Success      201  {object}  response.ResponseData
// @Failure      400  {object}  response.ErrorResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /chat/createRoom [post]
// @Security     BearerAuth
func (c *cChat) CreateRoom(ctx *gin.Context) {
	var params model.CreateRoom
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}
	codeRs, err := service.NewIChat().CreateRoom(ctx, &params)
	if err != nil {
		response.ErrorResponse(ctx, codeRs, err.Error())
	}
	response.SuccessResponse(ctx, codeRs, nil)
}
