package usecase

import (
	"av-merch-shop/internal/entities"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetInfo_AllErrors(t *testing.T) {
	ctx := t.Context()

	t.Run("get inventory error", func(t *testing.T) {
		ur := new(MockUserRepo)
		tr := new(MockTransactionRepo)
		or := new(MockOrderRepo)
		cache := new(MockUserInfoCache)

		ur.On("GetByID", ctx, 1).Return(&entities.User{ID: 1}, nil)
		expectedErr := errors.New("inventory error")
		or.On("GetUserInventory", ctx, 1).Return(nil, expectedErr)
		cache.On("GetUserInfo", ctx, 1).Return(nil, errors.New("cache miss"))

		tr.On("GetUserBalance", ctx, 1).Return(0, nil)
		tr.On("GetOutgoingForUser", ctx, 1).Return(&entities.UserSent{}, nil)
		tr.On("GetIncomingForUser", ctx, 1).Return(&entities.UserReceived{}, nil)

		uc := NewInfoUseCase(ur, tr, or, cache)
		info, err := uc.GetInfo(ctx, 1)

		assert.Error(t, err)
		assert.Nil(t, info)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("cache hit returns early", func(t *testing.T) {
		ur := new(MockUserRepo)
		tr := new(MockTransactionRepo)
		or := new(MockOrderRepo)
		cache := new(MockUserInfoCache)

		cachedInfo := &entities.UserInfo{
			Coins:     100,
			Inventory: entities.UserInventory{},
			CoinHistory: entities.UserCoinHistory{
				Received: entities.UserReceived{},
				Sent:     entities.UserSent{},
			},
		}

		cache.On("GetUserInfo", ctx, 1).Return(cachedInfo, nil)

		uc := NewInfoUseCase(ur, tr, or, cache)
		info, err := uc.GetInfo(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, cachedInfo, info)

		ur.AssertNotCalled(t, "GetByID")
		tr.AssertNotCalled(t, "GetUserBalance")
		or.AssertNotCalled(t, "GetUserInventory")
	})

}

func TestGetInfo_Success(t *testing.T) {
	ctx := t.Context()
	ur := new(MockUserRepo)
	tr := new(MockTransactionRepo)
	or := new(MockOrderRepo)
	cache := new(MockUserInfoCache)

	user := &entities.User{ID: 1}
	inventory := &entities.UserInventory{}
	sent := &entities.UserSent{}
	received := &entities.UserReceived{}

	ur.On("GetByID", ctx, 1).Return(user, nil)
	or.On("GetUserInventory", ctx, 1).Return(inventory, nil)
	tr.On("GetUserBalance", ctx, 1).Return(100, nil)
	tr.On("GetOutgoingForUser", ctx, 1).Return(sent, nil)
	tr.On("GetIncomingForUser", ctx, 1).Return(received, nil)
	cache.On("GetUserInfo", ctx, 1).Return(nil, errors.New("cache miss"))
	cache.On("SetUserInfo", ctx, 1, mock.AnythingOfType("entities.UserInfo")).Return(nil)

	uc := NewInfoUseCase(ur, tr, or, cache)
	info, err := uc.GetInfo(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, 100, info.Coins)
}
