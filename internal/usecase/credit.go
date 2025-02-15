package usecase

import (
	"av-merch-shop/internal/entities"
	"context"

	"github.com/google/uuid"
)

type creditUseCase struct {
	transactionManager TransactionManager
	transactionRepo    TransactionRepo
	userRepo           UserRepo
	userInfoCache      UserInfoCache
}

type CreditUseCase interface {
	Credit(ctx context.Context, data *entities.CreditData) (*entities.CreditTransactionResult, error)
}

func NewCreditUseCase(tm TransactionManager, tr TransactionRepo, ur UserRepo, uic UserInfoCache) CreditUseCase {
	return &creditUseCase{
		transactionManager: tm,
		transactionRepo:    tr,
		userRepo:           ur,
		userInfoCache:      uic,
	}
}

func (u *creditUseCase) Credit(
	ctx context.Context,
	data *entities.CreditData,
) (
	*entities.CreditTransactionResult, error) {

	var tr *entities.Transaction
	var balance int
	var user *entities.User

	err := u.transactionManager.Do(ctx, func(ctx context.Context) error {
		// проверяем наличие пользователя, если его не существует, то вернём NotFound
		var err error
		user, err = u.userRepo.GetByUsername(ctx, data.Username)
		if err != nil {
			return err
		}
		// проверяем баланс пользователя
		balance, err = u.transactionRepo.GetUserBalance(ctx, user.ID)
		if err != nil {
			return err
		}

		// если после суммирования баланс отрицательный, то это транзакция списания
		// и если она попытается сделать баланс ниже нуля — меняем размер транзакции
		// на остаток баланса
		if balance+data.Amount < 0 {
			data.Amount = -balance
		}

		// записываем транзакцию в БД
		tr, err = u.transactionRepo.Create(ctx, entities.TransactionData{
			UserID:          user.ID,
			Amount:          data.Amount,
			TransactionType: entities.TransactionCredit,
			ReferenceID:     uuid.New(),
		})
		if err != nil {
			return err
		}

		u.userInfoCache.ExpireUserInfo(ctx, user.ID)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &entities.CreditTransactionResult{
		NewAmount:   balance + data.Amount,
		ReferenceID: tr.ReferenceID,
	}, nil
}
