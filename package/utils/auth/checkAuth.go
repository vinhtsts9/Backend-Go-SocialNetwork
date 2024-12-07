package auth

import (
	"context"
	"encoding/json"
	"go-ecommerce-backend-api/m/v2/global"
	model "go-ecommerce-backend-api/m/v2/internal/models"
	"go-ecommerce-backend-api/m/v2/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func CheckAuth(ctx *gin.Context, token string) *jwt.StandardClaims {
	tokenString, ok := ExtractBearerToken(token)
	if !ok {
		response.ErrorResponse(ctx, http.StatusUnauthorized, "Missing or invalid token")
	}
	claims, err := VerifyTokenSubject(tokenString)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusUnauthorized, "Missing or invalid token")
	}
	return claims
}

func GetUserIdFromToken(ctx *gin.Context, token string) int {
	claims := CheckAuth(ctx, token)

	Result, err := global.Rdb.Get(context.Background(), claims.Subject).Result()
	if err != nil {
		response.ErrorResponse(ctx, http.StatusUnauthorized, "Session expried or invalid")
	}
	var userInfo model.UserInfo
	err = json.Unmarshal([]byte(Result), &userInfo)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "Internal Server Error")
	}

	userId := userInfo.UserID
	return int(userId)
}
