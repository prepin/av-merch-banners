package usecase

import (
	"av-merch-shop/internal/entities"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCredit_Errors(t *testing.T) {
	ctx := t.Context()

	t.Run("get user balance error", func(t *testing.T) {
		tm := new(MockTransactionManager)
		tr := new(MockTransactionRepo)
		ur := new(MockUserRepo)

		expectedErr := errors.New("balance error")

		ur.On("GetByUsername", ctx, "user").Return(&entities.User{ID: 1}, nil)
		tr.On("GetUserBalance", ctx, 1).Return(0, expectedErr)

		tm.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).
			Run(func(args mock.Arguments) {
				f := args.Get(1).(func(context.Context) error)
				f(ctx)
			}).
			Return(expectedErr)

		uc := NewCreditUseCase(tm, tr, ur)
		result, err := uc.Credit(ctx, &entities.CreditData{Username: "user"})

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
	})
}

func TestCredit_CreateTransactionError(t *testing.T) {
	ctx := t.Context()
	tm := new(MockTransactionManager)
	tr := new(MockTransactionRepo)
	ur := new(MockUserRepo)

	user := &entities.User{ID: 1, Username: "user"}
	ur.On("GetByUsername", ctx, "user").Return(user, nil)
	tr.On("GetUserBalance", ctx, user.ID).Return(100, nil)

	expectedErr := errors.New("transaction creation error")
	tr.On("Create", ctx, mock.AnythingOfType("entities.TransactionData")).Return(nil, expectedErr)

	tm.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			f := args.Get(1).(func(context.Context) error)
			f(ctx)
		}).
		Return(expectedErr)

	uc := NewCreditUseCase(tm, tr, ur)
	result, err := uc.Credit(ctx, &entities.CreditData{
		Username: "user",
		Amount:   50,
	})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedErr, err)
	ur.AssertExpectations(t)
	tr.AssertExpectations(t)
}
