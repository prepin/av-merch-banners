package usecase

import (
	"av-merch-shop/internal/entities"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCredit_Errors(t *testing.T) {
	ctx := t.Context()

	t.Run("get user balance error", func(t *testing.T) {
		tm := new(MockTransactionManager)
		tr := new(MockTransactionRepo)
		ur := new(MockUserRepo)
		uic := new(MockUserInfoCache)

		expectedErr := errors.New("balance error")

		ur.On("GetByUsername", ctx, "user").Return(&entities.User{ID: 1}, nil)
		tr.On("GetUserBalance", ctx, 1).Return(0, expectedErr)

		tm.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).
			Run(func(args mock.Arguments) {
				f := args.Get(1).(func(context.Context) error)
				f(ctx)
			}).
			Return(expectedErr)

		uc := NewCreditUseCase(tm, tr, ur, uic)
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
	uic := new(MockUserInfoCache)

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

	uc := NewCreditUseCase(tm, tr, ur, uic)
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

func TestCredit_Success(t *testing.T) {
	ctx := t.Context()
	tm := new(MockTransactionManager)
	tr := new(MockTransactionRepo)
	ur := new(MockUserRepo)
	uic := new(MockUserInfoCache)

	user := &entities.User{ID: 1, Username: "user"}
	initialBalance := 100
	creditAmount := 50
	referenceID := uuid.New()

	ur.On("GetByUsername", ctx, "user").Return(user, nil)
	tr.On("GetUserBalance", ctx, user.ID).Return(initialBalance, nil)
	tr.On("Create", ctx, mock.MatchedBy(func(data entities.TransactionData) bool {
		return data.UserID == user.ID &&
			data.Amount == creditAmount &&
			data.TransactionType == entities.TransactionCredit
	})).Return(&entities.Transaction{
		ReferenceID: referenceID,
	}, nil)

	uic.On("ExpireUserInfo", ctx, user.ID)

	tm.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			f := args.Get(1).(func(context.Context) error)
			f(ctx)
		}).
		Return(nil)

	uc := NewCreditUseCase(tm, tr, ur, uic)
	result, err := uc.Credit(ctx, &entities.CreditData{
		Username: "user",
		Amount:   creditAmount,
	})

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, initialBalance+creditAmount, result.NewAmount)
	assert.Equal(t, referenceID, result.ReferenceID)

	ur.AssertExpectations(t)
	tr.AssertExpectations(t)
	uic.AssertExpectations(t)
}
