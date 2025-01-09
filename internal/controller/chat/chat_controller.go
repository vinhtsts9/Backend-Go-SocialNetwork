package chat

import (
	model "go-ecommerce-backend-api/m/v2/internal/models"
	"go-ecommerce-backend-api/m/v2/internal/service"
	"go-ecommerce-backend-api/m/v2/package/utils/auth"
	"go-ecommerce-backend-api/m/v2/response"
	"strconv"

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

func (c *cChat) GetChatHistory(ctx *gin.Context) {
	roomID := ctx.Param("room_id")
	roomId, err := strconv.ParseInt(roomID, 10, 64)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeGetMessage, err.Error())
		return
	}
	codeRs, err, data := service.NewIChat().GetChatHistory(ctx, int(roomId))
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeGetMessage, err.Error())
		return
	}
	response.SuccessResponse(ctx, codeRs, data)

}
func (c *cChat) GetUserNickName(ctx *gin.Context) {
	codeRs, data, err := service.NewIChat().GetUserNickName(ctx)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeGetMessage, err.Error())
		return
	}
	Result := struct {
		UserNickName string `json:"user_nickname"`
	}{
		UserNickName: data,
	}
	response.SuccessResponse(ctx, codeRs, Result)

}
func (c *cChat) GetRoomByUserId(ctx *gin.Context) {

	userInfo := auth.GetUserInfoFromContext(ctx)
	codeRs, data, err := service.NewIChat().GetRoomChatByUserId(ctx, uint64(userInfo.UserID))
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeGetMessage, err.Error())
		return
	}
	response.SuccessResponse(ctx, codeRs, data)
}

func (c *cChat) DeleteMemberFromGroup(ctx *gin.Context) {
	userIdP := ctx.Param("user_id")
	userId, err := strconv.ParseUint(userIdP, 10, 64)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeGetMessage, err.Error())
		return
	}
	roomIdP := ctx.Param("room_id")
	roomId, err := strconv.ParseInt(roomIdP, 10, 64)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeGetMessage, err.Error())
		return
	}

	codeRs, Rs, err := service.NewIChat().DeleteMemberFromGroup(ctx, userId, roomId)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeGetMessage, err.Error())
		return
	}
	response.SuccessResponse(ctx, codeRs, Rs)
}

func (c *cChat) GetMemberGroup(ctx *gin.Context) {
	roomIdP := ctx.Param("room_id")
	roomId, err := strconv.ParseInt(roomIdP, 10, 64)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeGetMessage, err.Error())
		return
	}

	codeRs, Rs, err := service.NewIChat().GetMemberGroup(ctx, roomId)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeGetMessage, err.Error())
		return
	}
	response.SuccessResponse(ctx, codeRs, Rs)
}
