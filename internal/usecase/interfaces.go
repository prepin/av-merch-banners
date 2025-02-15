package usecase

import (
	"av-merch-shop/internal/entities"
	"context"
	"time"
)

type Repos struct {
	TransactionManager TransactionManager
	UserRepo           UserRepo
	TransactionRepo    TransactionRepo
	ItemRepo           ItemRepo
	OrderRepo          OrderRepo
	UserInfoCache      UserInfoCache
}

type Services struct {
	TokenService TokenService
	HashService  HashService
}

type TransactionManager interface {
	Do(ctx context.Context, f func(ctx context.Context) error) error
}

type UserRepo interface {
	GetByID(ctx context.Context, userID int) (*entities.User, error)
	GetByUsername(ctx context.Context, username string) (*entities.User, error)
	Create(ctx context.Context, data entities.UserData) (*entities.User, error)
}

type TransactionRepo interface {
	GetUserBalance(ctx context.Context, userID int) (int, error)
	GetIncomingForUser(ctx context.Context, userID int) (*entities.UserReceived, error)
	GetOutgoingForUser(ctx context.Context, userID int) (*entities.UserSent, error)
	Create(ctx context.Context, data entities.TransactionData) (*entities.Transaction, error)
}

type OrderRepo interface {
	Create(ctx context.Context, data entities.OrderData) (*entities.Order, error)
	GetUserInventory(ctx context.Context, userID int) (*entities.UserInventory, error)
}

type ItemRepo interface {
	GetItemByName(ctx context.Context, itemName string) (*entities.Item, error)
}

type UserInfoCache interface {
	GetUserInfo(ctx context.Context, userID int) (*entities.UserInfo, error)
	SetUserInfo(ctx context.Context, userID int, info entities.UserInfo) error
	ExpireUserInfo(ctx context.Context, userID int)
}

type TokenService interface {
	GenerateToken(userID int, username string, role string) (string, error)
}

type HashService interface {
	HashPassword(password string) (string, error)
	CompareWithPassword(hashed string, password string) bool
}

type CacheService interface {
	Get(ctx context.Context, key string, dest any) error
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
}
