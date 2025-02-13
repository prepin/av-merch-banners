package usecase

import (
	"av-merch-shop/internal/entities"
	"av-merch-shop/internal/errs"
	"context"
	"errors"

	"github.com/google/uuid"
)

type AuthUseCase struct {
	transactionManager TransactionManager
	userRepo           UserRepo
	transactionRepo    TransactionRepo
	tokenService       TokenService
	hashService        HashService
}

const DefaultBalanceForUser = 1000

func NewAuthUsecase(tm TransactionManager, ur UserRepo, tr TransactionRepo, ts TokenService, hs HashService) *AuthUseCase {
	return &AuthUseCase{
		transactionManager: tm,
		userRepo:           ur,
		transactionRepo:    tr,
		tokenService:       ts,
		hashService:        hs,
	}
}

func (u *AuthUseCase) SignIn(ctx context.Context, username string, password string) (string, error) {
	var token string

	err := u.transactionManager.Do(ctx, func(ctx context.Context) error {
		user, err := u.userRepo.GetByUsername(ctx, username)

		switch {

		// юзера в базе не нашли, создаём юзера
		case err != nil && errors.Is(err, errs.ErrNotFoundError):
			token, err = u.getTokenForNewUser(ctx, username, password)
			if err != nil {
				return err
			}
			return nil

		// что-то пошло не так
		case err != nil:
			return err

		// юзер нашёлся, генерим токен
		default:
			token, err = u.getTokenForExistingUser(user, password)
			if err != nil {
				return err
			}
			return nil
		}
	})

	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *AuthUseCase) getTokenForNewUser(ctx context.Context, username, password string) (string, error) {
	hashed, err := u.hashService.HashPassword(password)
	if err != nil {
		return "", err
	}

	// создаём пользователя
	newUser, err := u.createUser(ctx, username, hashed)
	if err != nil {
		return "", err
	}

	// начисляем монетки
	err = u.creditInitialAmount(ctx, newUser.ID)
	if err != nil {
		return "", err
	}

	// возвращаем токен
	return u.tokenService.GenerateToken(
		newUser.ID,
		newUser.Username,
		newUser.Role,
	)
}

func (u *AuthUseCase) createUser(ctx context.Context, username, hashedPassword string) (*entities.User, error) {
	return u.userRepo.Create(ctx, entities.UserData{
		Username:       username,
		HashedPassword: hashedPassword,
		Role:           entities.RoleUser,
	})
}

func (u *AuthUseCase) creditInitialAmount(ctx context.Context, userID int) error {
	_, err := u.transactionRepo.Create(ctx, entities.TransactionData{
		UserID:          userID,
		Amount:          DefaultBalanceForUser,
		TransactionType: entities.TransactionCredit,
		ReferenceId:     uuid.New(),
	})
	return err
}

func (u *AuthUseCase) getTokenForExistingUser(user *entities.User, password string) (string, error) {
	if !u.hashService.CompareWithPassword(user.HashedPassword, password) {
		return "", errs.ErrNoAccessError
	}

	return u.tokenService.GenerateToken(
		user.ID,
		user.Username,
		user.Role,
	)
}
