package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/tuantran0910/rainbow/config"
)

func GenerateToken(userId uuid.UUID) (string, error) {
	// Get the secret key from the config
	cfg, err := config.GetConfig()
	if err != nil {
		return "", err
	}
	jwtSecret := []byte(cfg.ServerConfig.JWTSecret)

	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenString string) (*jwt.Token, error) {
	// Get the secret key from the config
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	jwtSecret := []byte(cfg.ServerConfig.JWTSecret)

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
}
