package usecase

import (
	"av-merch-shop/internal/entities"
	"context"

	"github.com/google/uuid"
)

type CreditUseCase struct {
	transactionManager TransactionManager
	transactionRepo    TransactionRepo
	userRepo           UserRepo
}

func NewCreditUseCase(tm TransactionManager, tr TransactionRepo, ur UserRepo) *CreditUseCase {
	return &CreditUseCase{
		transactionManager: tm,
		transactionRepo:    tr,
		userRepo:           ur,
	}
}

func (u *CreditUseCase) Credit(
	ctx context.Context,
	data *entities.CreditData,
) (
	*entities.CreditTransactionResult, error) {

	var tr *entities.Transaction
	var balance int64
	var user *entities.User

	err := u.transactionManager.Do(ctx, func(ctx context.Context) error {
		// проверяем наличие пользователя, если его не существует, то вернём NotFound
		var err error
		user, err = u.userRepo.GetByUsername(ctx, data.Username)
		if err != nil {
			return err
		}
		// проверяем баланс пользователя
		balanceInt, err := u.transactionRepo.GetUserBalance(ctx, user.ID)
		if err != nil {
			return err
		}
		balance = int64(balanceInt)

		// если после суммирования баланс отрицательный, то это транзакция списания
		// и она попытается сделать баланс ниже нуля — меняем размер транзакции
		// на остаток баланса
		if balance+int64(data.Amount) < 0 {
			data.Amount = int(-balance)
		}

		// записываем транзакцию в БД
		tr, err = u.transactionRepo.CreateTransaction(ctx, entities.TransactionData{
			UserID:          user.ID,
			Amount:          data.Amount,
			TransactionType: entities.TransactionCredit,
			ReferenceId:     uuid.New(),
		})
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &entities.CreditTransactionResult{
		NewAmount:   int(balance + int64(data.Amount)),
		ReferenceID: tr.ReferenceId,
	}, nil
}
