package auth

import (
	"av-merch-shop/config"
	"log/slog"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int
	Role   string
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
