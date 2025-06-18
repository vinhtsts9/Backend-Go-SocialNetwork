package auth

import (
	"encoding/json"
	"fmt"
	"go-ecommerce-backend-api/m/v2/global"
	model "go-ecommerce-backend-api/m/v2/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func CheckAuth(token string) error {
	tokenString, ok := ExtractBearerToken(token)
	if !ok {
		return fmt.Errorf("don't extract bearer token")
	}
	_, err := VerifyTokenSubject(tokenString)
	if err != nil {
		return fmt.Errorf("Check auth failed", err)
	}
	return nil
}
func CheckAuthForWebsocket(token string) (*model.UserInfo, error) {
	tokenString, ok := ExtractBearerToken(token)
	if !ok {
		return &model.UserInfo{}, fmt.Errorf("don't extract bearer token")
	}
	claims, err := VerifyTokenSubject(tokenString)
	if err != nil {
		return &model.UserInfo{}, fmt.Errorf("Check auth failed", err)
	}

	subjectStr := claims.Subject
	if subjectStr == "" {
		global.Logger.Sugar().Error("Subject is empty in JWT claims")
		return &model.UserInfo{}, err
	}

	var userInfo model.UserInfo
	err = json.Unmarshal([]byte(subjectStr), &userInfo)
	if err != nil {
		global.Logger.Sugar().Errorf("JSON unmarshal error: %v", err)
		return &model.UserInfo{}, err
	}

	return &userInfo, nil
}
func GetUserInfoFromContext(ctx *gin.Context) *model.UserInfo {
	claimsValue := ctx.Request.Context().Value("claims")
	claims, ok := claimsValue.(*jwt.StandardClaims)
	if !ok {
		// Xử lý lỗi khi ép kiểu không thành công
		global.Logger.Sugar().Error("Claims are not of type jwt.StandardClaims")
		return &model.UserInfo{}
	}

	subjectStr := claims.Subject
	if subjectStr == "" {
		global.Logger.Sugar().Error("Subject is empty in JWT claims")
		return &model.UserInfo{}
	}

	var userInfo model.UserInfo
	err := json.Unmarshal([]byte(subjectStr), &userInfo)
	if err != nil {
		global.Logger.Sugar().Errorf("JSON unmarshal error: %v", err)
		return &model.UserInfo{}
	}

	return &userInfo
}
