package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tuantran0910/rainbow/config"
)

const (
	TokenExpiration = time.Minute * 60 // 1 hours
	TokenIssuer     = "rainbow-api"
)

func GenerateToken(userId uuid.UUID) (string, error) {
	// Get the secret key from the config
	cfg, err := config.GetConfig()
	if err != nil {
		return "", err
	}
	jwtSecret := []byte(cfg.ServerConfig.JWTSecret)

	now := time.Now()
	claims := jwt.MapClaims{
		"user_id": userId.String(),
		"exp":     now.Add(TokenExpiration).Unix(),
		"iat":     now.Unix(),          // issued at
		"nbf":     now.Unix(),          // not before
		"iss":     TokenIssuer,         // issuer
		"jti":     uuid.New().String(), // JWT ID for uniqueness
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	// Get the secret key from the config
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	jwtSecret := []byte(cfg.ServerConfig.JWTSecret)

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
}

func ValidateTokenClaims(token *jwt.Token) error {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("invalid token claims")
	}

	if iss, ok := claims["iss"].(string); !ok || iss != TokenIssuer {
		return fmt.Errorf("invalid token issuer")
	}

	if userIDStr, ok := claims["user_id"].(string); !ok {
		return fmt.Errorf("user_id missing in token")
	} else if _, err := uuid.Parse(userIDStr); err != nil {
		return fmt.Errorf("invalid user_id format in token")
	}

	return nil
}
