package account

import (
	"database/sql"
	"fmt"
	"go-ecommerce-backend-api/m/v2/global"
	"go-ecommerce-backend-api/m/v2/internal/database"
	model "go-ecommerce-backend-api/m/v2/internal/models"
	"go-ecommerce-backend-api/m/v2/internal/service"
	"go-ecommerce-backend-api/m/v2/package/utils/auth"
	"go-ecommerce-backend-api/m/v2/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// management login user
var Login = new(cUserLogin)

type cUserLogin struct{}

// UpdatePasswordRegister
// @Summary      UpdatePasswordRegister
// @Description  UpdatePasswordRegister
// @Tags         accounts management
// @Accept       json
// @Produce      json
// @Param        payload body model.UpdatePasswordRegisterInput true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /user/update_pass_register [post]
func (c *cUserLogin) UpdatePasswordRegister(ctx *gin.Context) {
	var params model.UpdatePasswordRegisterInput
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}
	result, err := service.UserLogin().UpdatePasswordRegister(ctx, &params)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}
	response.SuccessResponse(ctx, response.ErrCodeSuccess, result)
}

func (c *cUserLogin) UpdateAvatar(ctx *gin.Context) {
	err := ctx.Request.ParseMultipartForm(10 << 20) // Giới hạn 10MB
	if err != nil {
		fmt.Println("DEBUG: Multipart Form Parse Error:", err)
		ctx.JSON(400, gin.H{"error": "Invalid form data"})
		return
	}
	fmt.Println("DEBUG: Multipart Form Parsed Successfully")

	userInfo := auth.GetUserInfoFromContext(ctx)

	file, _, err := ctx.Request.FormFile("avatar")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không tìm thấy file"})
		return
	}
	defer file.Close()

	uploadResp, err := global.Cloudinary.UploadImageToCloudinaryFromReader(file, "avatar")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Upload thất bại"})
		return
	}
	params := database.UpdateAvatarParams{
		UserAvatar: sql.NullString{
			String: uploadResp,
			Valid:  true,
		},
		UserID: uint64(userInfo.UserID),
	}
	_, err = global.Store.UpdateAvatar(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Không cập nhật được database"})
		return
	}
	result := model.UpdateAvatar{
		UserAvatar: uploadResp,
	}
	response.SuccessResponse(ctx, response.ErrCodeSuccess, result)
}

// Verify OTP Login By User
// @Summary      Verify OTP Login By User
// @Description  Verify OTP Login By User
// @Tags         accounts management
// @Accept       json
// @Produce      json
// @Param        payload body model.VerifyInput true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /user/verify_account [post]
func (c *cUserLogin) VerifyOTP(ctx *gin.Context) {
	var params model.VerifyInput
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}
	result, err := service.UserLogin().VerifyOTP(ctx, &params)
	if err != nil {
		response.ErrorResponse(ctx, response.ErrInvalidOTP, err.Error())
		return
	}
	response.SuccessResponse(ctx, response.ErrCodeSuccess, result)
}

// User Login
// @Summary      User Login
// @Description  User Login
// @Tags         accounts management
// @Accept       json
// @Produce      json
// @Param        payload body model.LoginInput true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /user/login [post]
func (c *cUserLogin) Login(ctx *gin.Context) {
	var params model.LoginInput
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}
	codeRs, dataRs, err := service.UserLogin().Login(ctx, &params)
	if err != nil {
		response.ErrorResponse(ctx, codeRs, err.Error())
		return
	}
	response.SuccessResponse(ctx, codeRs, dataRs)
}

// User Registration documentation
// @Summary      User Registration
// @Description  When user is registered send otp to email
// @Tags         accounts management
// @Accept       json
// @Produce      json
// @Param        payload body model.RegisterInput true "payload"
// @Success      200  {object}  response.ResponseData
// @Failure      500  {object}  response.ErrorResponseData
// @Router       /user/register [post]
func (c *cUserLogin) Register(ctx *gin.Context) {
	var params model.RegisterInput
	if err := ctx.ShouldBindJSON(&params); err != nil {
		response.ErrorResponse(ctx, response.ErrCodeParamInvalid, err.Error())
		return
	}
	codeStatus, err := service.UserLogin().Register(ctx, &params)
	if err != nil {
		global.Logger.Error("Error registing user OTP", zap.Error(err))
		response.ErrorResponse(ctx, codeStatus, err.Error())
		return
	}
	response.SuccessResponse(ctx, codeStatus, nil)
}
