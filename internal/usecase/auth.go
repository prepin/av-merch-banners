package usecase

import (
	"av-merch-shop/internal/errs"
	"errors"
)

type AuthUseCase struct {
	userRepo     UserRepo
	tokenService TokenService
}

func NewAuthUsecase(ur UserRepo, ts TokenService) *AuthUseCase {
	return &AuthUseCase{
		userRepo:     ur,
		tokenService: ts,
	}
}

func (u *AuthUseCase) SignIn(username string, password string) (string, error) {
	// проверить, есть ли пользователь в базе
	user, err := u.userRepo.GetByUsername(username)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound{}) {
			// если нет, то создать пользователя и вернуть для него токен
			u.userRepo.Create()
			token, err := u.tokenService.GenerateToken(-1, user.Username, user.Role)
			if err != nil {
				return "", err
			}
			return token, nil
		}
		// что-то пошло не так, отдаём внутреннюю ошибку
		return "", err
	}
	// если есть, то проверить пароль на валидность
	if user.CompareWithPassword(password) {
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
