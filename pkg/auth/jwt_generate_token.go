package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (js *JWTService) GenerateToken(userID int, username string, role string) (string, error) {
	if string(js.secretKey) == "default" {
		js.logger.Warn("JWT Secret key not set! Using default value. Provide secret key in AV_SECRET environment variable.")
	}

	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(js.secretKey)
}
