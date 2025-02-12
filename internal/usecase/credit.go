package usecase

import (
	"av-merch-shop/internal/entities"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type CreditUseCase struct {
	transactionRepo TransactionRepo
	userRepo        UserRepo
}

func NewCreditUseCase(tr TransactionRepo, ur UserRepo) *CreditUseCase {
	return &CreditUseCase{
		transactionRepo: tr,
		userRepo:        ur,
	}
}

func (u *CreditUseCase) Credit(
	ctx context.Context,
	data *entities.CreditData,
) (
	*entities.CreditTransactionResult,
	error) {
	// начинаем транзакцию в базе
	// TODO:

	// проверяем наличие пользователя, если его не существует, то вернём NotFound
	user, err := u.userRepo.GetByUsername(ctx, data.Username)
	if err != nil {
		// TODO: вернуть NotFound если ошибка именно такая
		return nil, err
	}
	fmt.Println("User", user)
	// проверяем баланс пользователя
	balance, err := u.transactionRepo.GetUserBalance(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	// если после суммирования баланс отрицательный, то это транзакция списания
	// и она попытается сделать баланс ниже нуля — меняем размер транзакции
	// на остаток баланса
	if balance+data.Amount < 0 {
		data.Amount = -balance
	}

	// записываем транзакцию в БД
	tr, err := u.transactionRepo.CreateTransaction(ctx, entities.TransactionData{
		UserID:          user.ID,
		Amount:          data.Amount,
		TransactionType: entities.TransactionCredit,
		ReferenceId:     uuid.New(),
	})
	if err != nil {
		return nil, err
	}

	// завершаем транзакцию в базе

	return &entities.CreditTransactionResult{
		NewAmount:   balance + data.Amount,
		ReferenceID: tr.ReferenceId,
	}, nil
}
