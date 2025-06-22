package middlewares

import (
	"strings"
	"time"

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

		tokenString := parts[1]
		if len(tokenString) == 0 {
			c.AbortWithStatusJSON(401, gin.H{"error": "Token cannot be empty"})
			return
		}

		// Parse the JWT token
		token, err := jwtHelper.ParseToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token format"})
			return
		}

		// Validate token
		if !token.Valid {
			c.AbortWithStatusJSON(401, gin.H{"error": "Token is invalid or expired"})
			return
		}

		// Extract the token claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token claims"})
			return
		}

		// Validate expiration explicitly (additional check)
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				c.AbortWithStatusJSON(401, gin.H{"error": "Token has expired"})
				return
			}
		} else {
			c.AbortWithStatusJSON(401, gin.H{"error": "Token expiration claim missing"})
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

		// Set user in the context for downstream handlers
		c.Set("user_id", userId)
		c.Next()
	}
}
