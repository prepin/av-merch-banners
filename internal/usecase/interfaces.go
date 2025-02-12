package usecase

import (
	"av-merch-shop/internal/entities"
	"context"
)

type Repos struct {
	TransactionManager TransactionManager
	UserRepo           UserRepo
	TransactionRepo    TransactionRepo
}

type Services struct {
	TokenService TokenService
	HashService  HashService
}

type TransactionManager interface {
	Do(ctx context.Context, f func(ctx context.Context) error) error
}

type UserRepo interface {
	GetByUsername(ctx context.Context, username string) (*entities.User, error)
	Create(ctx context.Context, data entities.UserData) (*entities.User, error)
}

type TransactionRepo interface {
	GetUserBalance(ctx context.Context, userId int) (int, error)
	CreateTransaction(ctx context.Context, data entities.TransactionData) (*entities.Transaction, error)
}

type TokenService interface {
	GenerateToken(userID int, username string, role string) (string, error)
}

type HashService interface {
	HashPassword(password string) (string, error)
	CompareWithPassword(hashed string, password string) bool
}
