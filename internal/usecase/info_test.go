package usecase

import (
	"av-merch-shop/internal/entities"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInfo_AllErrors(t *testing.T) {
	ctx := context.Background()

	t.Run("get inventory error", func(t *testing.T) {
		ur := new(MockUserRepo)
		tr := new(MockTransactionRepo)
		or := new(MockOrderRepo)

		ur.On("GetByID", ctx, 1).Return(&entities.User{ID: 1}, nil)
		expectedErr := errors.New("inventory error")
		or.On("GetUserInventory", ctx, 1).Return(nil, expectedErr)

		uc := NewInfoUseCase(ur, tr, or)
		info, err := uc.GetInfo(ctx, 1)

		assert.Error(t, err)
		assert.Nil(t, info)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("get balance error", func(t *testing.T) {
		ur := new(MockUserRepo)
		tr := new(MockTransactionRepo)
		or := new(MockOrderRepo)

		ur.On("GetByID", ctx, 1).Return(&entities.User{ID: 1}, nil)
		or.On("GetUserInventory", ctx, 1).Return(&entities.UserInventory{}, nil)
		expectedErr := errors.New("balance error")
		tr.On("GetUserBalance", ctx, 1).Return(0, expectedErr)

		uc := NewInfoUseCase(ur, tr, or)
		info, err := uc.GetInfo(ctx, 1)

		assert.Error(t, err)
		assert.Nil(t, info)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("get outgoing transactions error", func(t *testing.T) {
		ur := new(MockUserRepo)
		tr := new(MockTransactionRepo)
		or := new(MockOrderRepo)

		ur.On("GetByID", ctx, 1).Return(&entities.User{ID: 1}, nil)
		or.On("GetUserInventory", ctx, 1).Return(&entities.UserInventory{}, nil)
		tr.On("GetUserBalance", ctx, 1).Return(100, nil)
		expectedErr := errors.New("outgoing error")
		tr.On("GetOutgoingForUser", ctx, 1).Return(nil, expectedErr)

		uc := NewInfoUseCase(ur, tr, or)
		info, err := uc.GetInfo(ctx, 1)

		assert.Error(t, err)
		assert.Nil(t, info)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("get incoming transactions error", func(t *testing.T) {
		ur := new(MockUserRepo)
		tr := new(MockTransactionRepo)
		or := new(MockOrderRepo)

		ur.On("GetByID", ctx, 1).Return(&entities.User{ID: 1}, nil)
		or.On("GetUserInventory", ctx, 1).Return(&entities.UserInventory{}, nil)
		tr.On("GetUserBalance", ctx, 1).Return(100, nil)
		tr.On("GetOutgoingForUser", ctx, 1).Return(&entities.UserSent{}, nil)
		expectedErr := errors.New("incoming error")
		tr.On("GetIncomingForUser", ctx, 1).Return(nil, expectedErr)

		uc := NewInfoUseCase(ur, tr, or)
		info, err := uc.GetInfo(ctx, 1)

		assert.Error(t, err)
		assert.Nil(t, info)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("get user error", func(t *testing.T) {
		ur := new(MockUserRepo)
		tr := new(MockTransactionRepo)
		or := new(MockOrderRepo)

		expectedErr := errors.New("user not found")
		ur.On("GetByID", ctx, 1).Return(nil, expectedErr)

		uc := NewInfoUseCase(ur, tr, or)
		info, err := uc.GetInfo(ctx, 1)

		assert.Error(t, err)
		assert.Nil(t, info)
		assert.Equal(t, expectedErr, err)
		ur.AssertExpectations(t)
	})
}
