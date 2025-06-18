package account

import (
	"go-ecommerce-backend-api/m/v2/package/utils/auth"

	"github.com/gin-gonic/gin"
)

func RefreshTokenHandler(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	userInfo, err := auth.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	newAccessToken, newRefreshToken, err := auth.RotateRefreshToken(req.RefreshToken, userInfo)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate tokens"})
		return
	}

	c.JSON(200, gin.H{
		"token":         newAccessToken,
		"refresh_token": newRefreshToken,
	})
}
