package usecase

import (
	"av-merch-shop/internal/entities"
	"av-merch-shop/internal/errs"
	"context"

	"github.com/google/uuid"
)

type SendCoinUseCase struct {
	transactionManager TransactionManager
	transactionRepo    TransactionRepo
	userRepo           UserRepo
}

func NewSendCoinUseCase(tm TransactionManager, tr TransactionRepo, ur UserRepo) *SendCoinUseCase {
	return &SendCoinUseCase{
		transactionManager: tm,
		transactionRepo:    tr,
		userRepo:           ur,
	}
}

func (u *SendCoinUseCase) Send(ctx context.Context, data *entities.TransferData) error {

	// не разрешаем перевести деньги в обратную сторону
	if data.Amount <= 0 {
		return errs.ErrIncorrectAmountError
	}

	err := u.transactionManager.Do(ctx, func(ctx context.Context) error {
		// проверить что существует отправитель
		sender, err := u.userRepo.GetByID(ctx, data.SenderID)
		if err != nil {
			return err
		}
		// проверить что существует получатель
		recipient, err := u.userRepo.GetByUsername(ctx, data.Recipient)
		if err != nil {
			return err
		}
		// проверить что у получателя достаточный баланс
		balance, err := u.transactionRepo.GetUserBalance(ctx, sender.ID)
		if err != nil {
			return err
		}
		if balance < data.Amount {
			return errs.ErrInsufficientFundsError
		}

		// создать две записи в транзакцию
		ref := uuid.New()
		outTrData := entities.TransactionData{
			UserID:         sender.ID,
			CounterpartyID: recipient.ID,
			Amount:         -data.Amount,
			ReferenceId:    ref,
		}

		inTrData := entities.TransactionData{
			UserID:          recipient.ID,
			CounterpartyID:  sender.ID,
			Amount:          data.Amount,
			TransactionType: entities.TransactionTransfer,
			ReferenceId:     ref,
		}

		_, err = u.transactionRepo.Create(ctx, outTrData)
		if err != nil {
			return err
		}

		_, err = u.transactionRepo.Create(ctx, inTrData)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
