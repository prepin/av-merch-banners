package usecase

import (
	"av-merch-shop/internal/entities"
	"context"
	"log"
	"sync"
)

type infoUseCase struct {
	transactionRepo TransactionRepo
	userRepo        UserRepo
	orderRepo       OrderRepo
	userInfoCache   UserInfoCache
}

type InfoUseCase interface {
	GetInfo(ctx context.Context, userID int) (*entities.UserInfo, error)
}

func NewInfoUseCase(ur UserRepo, tr TransactionRepo, or OrderRepo, uic UserInfoCache) InfoUseCase {
	return &infoUseCase{
		transactionRepo: tr,
		userRepo:        ur,
		orderRepo:       or,
		userInfoCache:   uic,
	}
}

func (u *infoUseCase) GetInfo(ctx context.Context, userID int) (*entities.UserInfo, error) {
	var cachedInfo *entities.UserInfo
	cachedInfo, err := u.userInfoCache.GetUserInfo(ctx, userID)
	if err == nil {
		return cachedInfo, nil
	}

	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	wg.Add(4)

	var (
		inventory *entities.UserInventory
		balance   int
		sent      *entities.UserSent
		received  *entities.UserReceived
		errChan   = make(chan error, 4)
	)

	go func() {
		defer wg.Done()
		inv, err := u.orderRepo.GetUserInventory(ctx, user.ID)
		if err != nil {
			errChan <- err
			return
		}
		inventory = inv
	}()

	go func() {
		defer wg.Done()
		bal, err := u.transactionRepo.GetUserBalance(ctx, user.ID)
		if err != nil {
			errChan <- err
			return
		}
		balance = bal
	}()

	go func() {
		defer wg.Done()
		s, err := u.transactionRepo.GetOutgoingForUser(ctx, user.ID)
		if err != nil {
			errChan <- err
			return
		}
		sent = s
	}()

	go func() {
		defer wg.Done()
		r, err := u.transactionRepo.GetIncomingForUser(ctx, user.ID)
		if err != nil {
			errChan <- err
			return
		}
		received = r
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	info := entities.UserInfo{
		Coins:     balance,
		Inventory: *inventory,
		CoinHistory: entities.UserCoinHistory{
			Received: *received,
			Sent:     *sent,
		},
	}

	err = u.userInfoCache.SetUserInfo(ctx, userID, info)
	if err != nil {
		log.Println("Failed to set cache for user:", userID)
	}

	return &info, nil
}
