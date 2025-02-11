package usecase

import (
	"av-merch-shop/internal/entities"
	"context"
)

type Repos struct {
	UserRepo UserRepo
}

type Services struct {
	TokenService TokenService
	HashService  HashService
}

type UserRepo interface {
	GetByUsername(ctx context.Context, username string) (*entities.User, error)
	Create(ctx context.Context, data entities.UserData) (*entities.User, error)
}

type TokenService interface {
	GenerateToken(userID int, username string, role string) (string, error)
}

type HashService interface {
	HashPassword(password string) (string, error)
	CompareWithPassword(hashed string, password string) bool
}
