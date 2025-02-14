package auth

import (
	"av-merch-shop/config"
	"errors"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID   int
	Username string
	Role     string
	jwt.RegisteredClaims
}

type JWTService struct {
	logger    *slog.Logger
	secretKey []byte
}

func NewJWTService(cfg *config.Config) *JWTService {
	return &JWTService{
		secretKey: cfg.Auth.SecretKey,
		logger:    cfg.Logger,
	}
}

func (js *JWTService) GenerateToken(userID int, username, role string) (string, error) {
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

func (js *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString, &Claims{},
		func(_ *jwt.Token,
		) (any, error) {
			return js.secretKey, nil
		})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
