package usecase

import (
	"av-merch-shop/internal/entities"
	"context"

	"github.com/google/uuid"
)

type orderUseCase struct {
	transactionManager TransactionManager
	transactionRepo    TransactionRepo
	userRepo           UserRepo
	itemRepo           ItemRepo
	orderRepo          OrderRepo
}

type OrderUseCase interface {
	Buy(ctx context.Context, data *entities.OrderRequest) error
}

func NewOrderUseCase(
	tm TransactionManager,
	tr TransactionRepo,
	ur UserRepo,
	ir ItemRepo,
	or OrderRepo,
) OrderUseCase {
	return &orderUseCase{
		transactionManager: tm,
		transactionRepo:    tr,
		userRepo:           ur,
		itemRepo:           ir,
		orderRepo:          or,
	}
}

func (u *orderUseCase) Buy(ctx context.Context, data *entities.OrderRequest) error {

	err := u.transactionManager.Do(ctx, func(ctx context.Context) error {
		// проверяем что вещь существует
		item, err := u.itemRepo.GetItemByName(ctx, data.ItemName)
		if err != nil {
			return err
		}
		// создаём транзакцию
		trnsct, err := u.transactionRepo.Create(ctx, entities.TransactionData{
			UserID:          data.UserID,
			Amount:          item.Cost,
			TransactionType: entities.TransactionOrder,
			ReferenceID:     uuid.New(),
		})
		if err != nil {
			return err
		}

		// создаём заказ
		_, err = u.orderRepo.Create(ctx, entities.OrderData{
			UserID:        data.UserID,
			ItemID:        item.ID,
			TransactionID: trnsct.ID,
		})
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
