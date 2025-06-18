package auth

import (
	"context"
	"fmt"
	"go-ecommerce-backend-api/m/v2/global"
	model "go-ecommerce-backend-api/m/v2/internal/models"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type PayloadClaims struct {
	jwt.StandardClaims
	UserInfo model.UserInfo
}

func GenTokenJWT(payload jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(global.Config.JWT.API_SECRET_KEY))
}
func CreateToken(infoUser string) (string, error) {
	timeEx := global.Config.JWT.JWT_EXPIRATION
	if timeEx == "" {
		timeEx = "1h"
	}
	expiration, err := time.ParseDuration(timeEx)
	if err != nil {
		return "", err
	}
	now := time.Now()
	expiresAt := now.Add(expiration)
	return GenTokenJWT(&PayloadClaims{
		StandardClaims: jwt.StandardClaims{
			Id:        uuid.New().String(),
			ExpiresAt: expiresAt.Unix(),
			IssuedAt:  now.Unix(),
			Issuer:    "Vũ Thế Vinh",
			Subject:   infoUser,
		},
	})
}
func CreateRefreshToken(userInfo string) (string, error) {
	refreshToken := uuid.New().String()
	expiration, err := time.ParseDuration(global.Config.JWT.REFRESH_EXPIRATION)
	if err != nil {
		return "", fmt.Errorf("invalid refresh token expiration config: %v", err)
	}
	ctx := context.Background()
	key := fmt.Sprintf("refresh:%s", refreshToken)
	err = global.Rdb.Set(ctx, key, userInfo, expiration).Err()
	if err != nil {
		return "", fmt.Errorf("failed to store refresh token: %v", err)
	}

	return refreshToken, nil
}
func RevokeRefreshToken(refreshToken string) error {
	ctx := context.Background()
	key := fmt.Sprintf("refresh:%s", refreshToken)
	if err := global.Rdb.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to revoke refresh token: %v", err)
	}
	return nil
}

// RotateRefreshToken replaces the old refresh token with a new one
func RotateRefreshToken(oldToken string, userInfo string) (string, string, error) {
	// Revoke the old token
	if err := RevokeRefreshToken(oldToken); err != nil {
		return "", "", fmt.Errorf("failed to revoke old refresh token: %v", err)
	}

	// Generate a new refresh token
	return GenerateTokens(userInfo)
}

func ParseJwtTokenSubject(token string) (*jwt.StandardClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(jwtToken *jwt.Token) (interface{}, error) {
		return []byte(global.Config.JWT.API_SECRET_KEY), nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*jwt.StandardClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

func VerifyTokenSubject(token string) (*jwt.StandardClaims, error) {
	claims, err := ParseJwtTokenSubject(token)
	if err != nil {
		return &jwt.StandardClaims{}, fmt.Errorf("Error verify token ", err)
	}
	if err = claims.Valid(); err != nil {
		return &jwt.StandardClaims{}, fmt.Errorf("Error verify token ", err)

	}
	return claims, nil
}
func GenerateTokens(userInfo string) (accessToken, refreshToken string, err error) {
	accessToken, err = CreateToken(userInfo)
	if err != nil {
		return "", "", fmt.Errorf("failed to create access token: %v", err)
	}

	refreshToken, err = CreateRefreshToken(userInfo)
	if err != nil {
		return "", "", fmt.Errorf("failed to create refresh token: %v", err)
	}

	return accessToken, refreshToken, nil
}
func ValidateRefreshToken(refreshToken string) (string, error) {
	ctx := context.Background()
	key := fmt.Sprintf("refresh:%s", refreshToken)

	// Sử dụng Get và xử lý lỗi khi khóa không tồn tại
	userInfo, err := global.Rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("refresh token không tồn tại")
	} else if err != nil {
		return "", fmt.Errorf("lỗi khi truy xuất redis: %v", err)
	}

	return userInfo, nil
}
