package usecase

import (
	"av-merch-shop/internal/entities"
	"context"

	"github.com/google/uuid"
)

type OrderUseCase struct {
	transactionManager TransactionManager
	transactionRepo    TransactionRepo
	userRepo           UserRepo
	itemRepo           ItemRepo
	orderRepo          OrderRepo
}

func NewOrderUseCase(
	tm TransactionManager,
	tr TransactionRepo,
	ur UserRepo,
	ir ItemRepo,
	or OrderRepo,
) *OrderUseCase {
	return &OrderUseCase{
		transactionManager: tm,
		transactionRepo:    tr,
		userRepo:           ur,
		itemRepo:           ir,
		orderRepo:          or,
	}
}

func (u *OrderUseCase) Buy(ctx context.Context, data *entities.OrderRequest) error {

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
			ReferenceId:     uuid.New(),
		})
		if err != nil {
			return err
		}

		// создаём заказ
		_, err = u.orderRepo.Create(ctx, entities.OrderData{
			UserID:        data.UserID,
			ItemId:        item.ID,
			TransactionId: trnsct.ID,
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
