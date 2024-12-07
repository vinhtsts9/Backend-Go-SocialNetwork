package middlewares

import (
	"context"
	"go-ecommerce-backend-api/m/v2/package/utils/auth"
	"log"

	"github.com/gin-gonic/gin"
)

func AuthenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uri := c.Request.URL.Path
		log.Println("uri request", uri)
		authHeader := c.GetHeader("Authorization")
		// check headers authorization
		jwtToken, valid := auth.ExtractBearerToken(authHeader)
		if !valid {
			c.AbortWithStatusJSON(401, gin.H{"code": 400011, "error": "Unauthorized", "description": ""})
			return
		}
		// validate jwt token by subject
		claims, err := auth.VerifyTokenSubject(jwtToken)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"code": 400012, "error": "InvalidToken", "description": ""})
			return
		}
		// update claims to context
		log.Println("claims:: UUID::", claims.Subject)
		ctx := context.WithValue(c.Request.Context(), "subjectUUID", claims.Subject)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}

}
