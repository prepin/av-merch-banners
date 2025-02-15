package usecase

import (
	"av-merch-shop/internal/entities"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBuy_Errors(t *testing.T) {
	ctx := t.Context()

	t.Run("item repo error", func(t *testing.T) {
		tm := new(MockTransactionManager)
		tr := new(MockTransactionRepo)
		ur := new(MockUserRepo)
		ir := new(MockItemRepo)
		or := new(MockOrderRepo)
		uic := new(MockUserInfoCache)

		expectedErr := errors.New("db error")

		ir.On("GetItemByName", ctx, "item1").Return(nil, expectedErr)

		tm.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).
			Run(func(args mock.Arguments) {
				f := args.Get(1).(func(context.Context) error)
				f(ctx)
			}).
			Return(expectedErr)

		uc := NewOrderUseCase(tm, tr, ur, ir, or, uic)
		err := uc.Buy(ctx, &entities.OrderRequest{ItemName: "item1", UserID: 1})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		tm.AssertExpectations(t)
		ir.AssertExpectations(t)
	})

	t.Run("transaction creation error", func(t *testing.T) {
		tm := new(MockTransactionManager)
		tr := new(MockTransactionRepo)
		ur := new(MockUserRepo)
		ir := new(MockItemRepo)
		or := new(MockOrderRepo)
		uic := new(MockUserInfoCache)

		expectedErr := errors.New("transaction error")

		ir.On("GetItemByName", ctx, "item1").Return(&entities.Item{
			ID:       1,
			Codename: "item1",
			Cost:     100,
		}, nil)

		tr.On("Create", ctx, mock.AnythingOfType("entities.TransactionData")).
			Return(nil, expectedErr)

		tm.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).
			Run(func(args mock.Arguments) {
				f := args.Get(1).(func(context.Context) error)
				f(ctx)
			}).
			Return(expectedErr)

		uc := NewOrderUseCase(tm, tr, ur, ir, or, uic)
		err := uc.Buy(ctx, &entities.OrderRequest{ItemName: "item1", UserID: 1})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		tm.AssertExpectations(t)
		ir.AssertExpectations(t)
		tr.AssertExpectations(t)
	})

	t.Run("order creation error", func(t *testing.T) {
		tm := new(MockTransactionManager)
		tr := new(MockTransactionRepo)
		ur := new(MockUserRepo)
		ir := new(MockItemRepo)
		or := new(MockOrderRepo)
		uic := new(MockUserInfoCache)

		expectedErr := errors.New("order error")

		ir.On("GetItemByName", ctx, "item1").Return(&entities.Item{
			ID:       1,
			Codename: "item1",
			Cost:     100,
		}, nil)

		tr.On("Create", ctx, mock.AnythingOfType("entities.TransactionData")).
			Return(&entities.Transaction{
				ID:     1,
				UserID: 1,
				Amount: 100,
			}, nil)

		or.On("Create", ctx, mock.AnythingOfType("entities.OrderData")).
			Return(nil, expectedErr)

		tm.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).
			Run(func(args mock.Arguments) {
				f := args.Get(1).(func(context.Context) error)
				f(ctx)
			}).
			Return(expectedErr)

		uc := NewOrderUseCase(tm, tr, ur, ir, or, uic)
		err := uc.Buy(ctx, &entities.OrderRequest{ItemName: "item1", UserID: 1})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		tm.AssertExpectations(t)
		ir.AssertExpectations(t)
		tr.AssertExpectations(t)
		or.AssertExpectations(t)
	})
}
