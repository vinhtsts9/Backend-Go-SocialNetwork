package account

import (
	model "go-ecommerce-backend-api/m/v2/internal/models"
	"go-ecommerce-backend-api/m/v2/internal/service"
	Context "go-ecommerce-backend-api/m/v2/package/utils/context"
	"go-ecommerce-backend-api/m/v2/response"
	"log"

	"github.com/gin-gonic/gin"
)

var TwoFA = new(sUser2FA)

type sUser2FA struct{}

// User Setup Two Factor Authentication
// @Summary      User Setup Two Factor Authentication
// @Description  User Setup Two Factor Authentication
// @Tags         account 2fa
// @Accept       json
// @Produce      json
// @Param Authorization header string true "Authorization token"
// @Param        payload body model.SetupTwoFactorAuthInput true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /user/two-factor/setup [post]
func (c *sUser2FA) SetupTwoFactorAuth(ctx *gin.Context) {
	var params model.SetupTwoFactorAuthInput
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeTwoFactorSetupFailed, "Missing or Invaid params")
		return
	}
	// get userId from uuid
	userId, err := Context.GetUserIdFromUUID(ctx.Request.Context())
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeTwoFactorSetupFailed, "userId is not invalid")
		return
	}
	log.Println("UserId", userId)
	params.UserId = uint32(userId)
	codeResult, err := service.UserLogin().SetupTwoFactorAuth(ctx, &params)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeTwoFactorSetupFailed, err.Error())
		return
	}
	response.SuccessResponse(ctx, codeResult, nil)
}
