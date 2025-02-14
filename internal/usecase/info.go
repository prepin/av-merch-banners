package usecase

import (
	"av-merch-shop/internal/entities"
	"context"
)

type infoUseCase struct {
	transactionRepo TransactionRepo
	userRepo        UserRepo
	orderRepo       OrderRepo
}

type InfoUseCase interface {
	GetInfo(ctx context.Context, userID int) (*entities.UserInfo, error)
}

func NewInfoUseCase(ur UserRepo, tr TransactionRepo, or OrderRepo) InfoUseCase {
	return &infoUseCase{
		transactionRepo: tr,
		userRepo:        ur,
		orderRepo:       or,
	}
}

func (u *infoUseCase) GetInfo(ctx context.Context, userID int) (*entities.UserInfo, error) {

	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	inventory, err := u.orderRepo.GetUserInventory(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	balance, err := u.transactionRepo.GetUserBalance(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	sent, err := u.transactionRepo.GetOutgoingForUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	received, err := u.transactionRepo.GetIncomingForUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	info := entities.UserInfo{

		Coins:     balance,
		Inventory: *inventory,
		CoinHistory: entities.UserCoinHistory{
			Received: *received,
			Sent:     *sent,
		},
	}

	return &info, nil
}
