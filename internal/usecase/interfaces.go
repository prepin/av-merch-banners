package usecase

import "av-merch-shop/internal/entities"

type Repos struct {
	UserRepo UserRepo
}

type Services struct {
	TokenService TokenService
}

type UserRepo interface {
	GetByUsername(username string) (*entities.User, error)
	Create()
}

type TokenService interface {
	GenerateToken(userID int, username string, role string) (string, error)
}
