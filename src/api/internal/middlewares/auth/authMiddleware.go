package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	jwtHelper "github.com/tuantran0910/rainbow/pkg/utils/jwt"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the value of the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header required"})
			return
		}

		// Check if the token is valid
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid Authorization header format"})
			return
		}

		// Parse the JWT token
		token, err := jwtHelper.ParseToken(parts[1])
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Extract the token claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token claims"})
			return
		}

		userIDStr, ok := claims["user_id"].(string)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "User ID missing in token"})
			return
		}

		userId, err := uuid.Parse(userIDStr)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid user ID format"})
			return
		}

		// Set user in the context
		c.Set("user_id", userId)
		c.Next()
	}
}
