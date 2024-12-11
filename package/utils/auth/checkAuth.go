package auth

import (
	"context"
	"encoding/json"
	"go-ecommerce-backend-api/m/v2/global"
	model "go-ecommerce-backend-api/m/v2/internal/models"
	"net/http"

	"github.com/golang-jwt/jwt"
)

func CheckAuth(w http.ResponseWriter, token string) *jwt.StandardClaims {
	tokenString, ok := ExtractBearerToken(token)
	if !ok {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return nil
	}
	claims, err := VerifyTokenSubject(tokenString)
	if err != nil {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return nil
	}
	return claims
}

func GetUserIdFromToken(w http.ResponseWriter, token string) int {

	claims := CheckAuth(w, token)
	if claims == nil {
		return 1
	}

	Result, err := global.Rdb.Get(context.Background(), claims.Subject).Result()
	if err != nil {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return 0

	}
	global.Logger.Sugar().Info(Result)

	var userInfo model.UserInfo
	err = json.Unmarshal([]byte(Result), &userInfo)
	if err != nil {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return 0
	}

	userId := userInfo.UserID
	return int(userId)
}
