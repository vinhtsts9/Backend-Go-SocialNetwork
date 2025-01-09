package auth

import (
	"context"
	"encoding/json"
	"go-ecommerce-backend-api/m/v2/global"
	model "go-ecommerce-backend-api/m/v2/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func CheckAuth(token string) *jwt.StandardClaims {
	tokenString, ok := ExtractBearerToken(token)
	if !ok {
		return nil
	}
	claims, err := VerifyTokenSubject(tokenString)
	if err != nil {
		return nil
	}
	return claims
}

func GetUserInfoFromToken(token string) model.UserInfo {

	claims := CheckAuth(token)
	if claims == nil {
		return model.UserInfo{}
	}

	Result, err := global.Rdb.Get(context.Background(), claims.Subject).Result()
	if err != nil {
		return model.UserInfo{}

	}
	var userInfo model.UserInfo
	err = json.Unmarshal([]byte(Result), &userInfo)
	if err != nil {
		return model.UserInfo{}
	}

	return userInfo
}
func GetUserInfoFromContext(ctx *gin.Context) model.UserInfo {

	Subject := ctx.Request.Context().Value("subjectUUID")
	subjectStr, ok := Subject.(string)
	if !ok {
		// Xử lý lỗi khi ép kiểu không thành công
		global.Logger.Sugar().Error("Subject is not a string")
		return model.UserInfo{}
	}
	Result, err := global.Rdb.Get(ctx, subjectStr).Result()
	if err != nil {
		return model.UserInfo{}
	}
	var userInfo model.UserInfo
	err = json.Unmarshal([]byte(Result), &userInfo)
	if err != nil {
		return model.UserInfo{}
	}
	return userInfo
}
