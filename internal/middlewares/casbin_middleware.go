package middlewares

import (
	"context"
	"encoding/json"
	"go-ecommerce-backend-api/m/v2/global"
	model "go-ecommerce-backend-api/m/v2/internal/models"
	"go-ecommerce-backend-api/m/v2/internal/service"
	"go-ecommerce-backend-api/m/v2/package/utils/auth"
	"go-ecommerce-backend-api/m/v2/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CasbinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		tokenString, ok := auth.ExtractBearerToken(authHeader)
		if !ok {
			response.ErrorResponse(c, http.StatusUnauthorized, "Missin or invalid token")
		}
		claims, err := auth.VerifyTokenSubject(tokenString)
		if err != nil {
			global.Logger.Sugar().Errorf("Invalid token %v", err)
			response.ErrorResponse(c, http.StatusUnauthorized, "Invalid Token")
			return
		}
		Result, err := global.Rdb.Get(context.Background(), claims.Subject).Result()
		if err != nil {
			response.ErrorResponse(c, http.StatusUnauthorized, "Session expried or invalid")
		}
		var userInfo model.UserInfo
		err = json.Unmarshal([]byte(Result), &userInfo)
		if err != nil {
			response.ErrorResponse(c, http.StatusInternalServerError, "Internal Server Error")
		}

		userId := userInfo.UserID
		sub := userId
		obj := c.Request.URL.Path
		act := c.Request.Method

		ok, err = global.Casbin.Enforce(sub, obj, act)
		if err != nil {
			global.Logger.Sugar().Error("Error enforcing policy: %v", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error 1"})
			c.Abort()
			return
		}
		if !ok {
			global.Logger.Sugar().Infof("Policy not found in cache, querying database for user: %s", sub)

			permissions, dbErr := service.RbacService().GetPermissionsByUserID(c, int64(sub))
			if dbErr != nil || len(permissions) == 0 {
				global.Logger.Sugar().Warnf("No permissions found for user: %s", sub)
				c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
				c.Abort()
				return
			}

			// Thêm quyền mới vào Casbin enforcer
			for _, perm := range permissions {

				global.Casbin.AddPolicy(sub, perm.Resource, perm.Action)

			}

			// Ghi lại chính sách vào file policy.csv
			if err := global.Casbin.SavePolicy(); err != nil {
				global.Logger.Sugar().Errorf("Failed to save policy: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				c.Abort()
				return
			}

			// Kiểm tra lại quyền sau khi cập nhật chính sách
			ok, err = global.Casbin.Enforce(sub, obj, act)
			if err != nil || !ok {
				global.Logger.Sugar().Warnf("Access denied for user: %s", sub)
				c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
