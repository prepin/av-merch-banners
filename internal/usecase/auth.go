package usecase

import (
	"av-merch-shop/internal/entities"
	"av-merch-shop/internal/errs"
	"context"
	"errors"
	"fmt"
)

type AuthUseCase struct {
	userRepo     UserRepo
	tokenService TokenService
	hashService  HashService
}

func NewAuthUsecase(ur UserRepo, ts TokenService, hs HashService) *AuthUseCase {
	return &AuthUseCase{
		userRepo:     ur,
		tokenService: ts,
		hashService:  hs,
	}
}

func (u *AuthUseCase) SignIn(ctx context.Context, username string, password string) (string, error) {
	// проверить, есть ли пользователь в базе
	user, err := u.userRepo.GetByUsername(ctx, username)
	fmt.Println(user, err)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound{}) {
			// если нет, то создать пользователя и вернуть для него токен
			hashed, err := u.hashService.HashPassword(password)
			if err != nil {
				return "", err
			}

			newUser, err := u.userRepo.Create(ctx, entities.UserData{
				Username:       username,
				HashedPassword: hashed,
				Role:           entities.RoleUser,
			})
			if err != nil {
				return "", err
			}

			token, err := u.tokenService.GenerateToken(
				newUser.ID,
				newUser.Username,
				newUser.Role,
			)
			if err != nil {
				return "", err
			}
			return token, nil
		}
		// что-то пошло не так, отдаём внутреннюю ошибку
		return "", err
	}
	// если есть, то проверить пароль на валидность
	if u.hashService.CompareWithPassword(user.HashedPassword, password) {
		// и если он валиден то отдать токен
		token, err := u.tokenService.GenerateToken(user.ID, user.Username, user.Role)
		if err != nil {
			return "", err
		}
		return token, nil
	}
	// а если невалиден, то отдаём ошибку доступа
	return "", errs.ErrNoAccessError
}
