package usecase

import (
	"av-merch-shop/internal/entities"
	"av-merch-shop/internal/errs"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSend_Errors(t *testing.T) {
	ctx := t.Context()

	t.Run("negative amount", func(t *testing.T) {
		tm := new(MockTransactionManager)
		tr := new(MockTransactionRepo)
		ur := new(MockUserRepo)

		uc := NewSendCoinUseCase(tm, tr, ur)
		err := uc.Send(ctx, &entities.TransferData{Amount: -1})

		assert.ErrorIs(t, err, errs.ErrIncorrectAmountError)
	})

	t.Run("insufficient funds", func(t *testing.T) {
		tm := new(MockTransactionManager)
		tr := new(MockTransactionRepo)
		ur := new(MockUserRepo)

		ur.On("GetByID", ctx, 1).Return(&entities.User{ID: 1}, nil)
		ur.On("GetByUsername", ctx, "recipient").Return(&entities.User{ID: 2}, nil)
		tr.On("GetUserBalance", ctx, 1).Return(50, nil) // Balance less than transfer amount

		tm.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).
			Run(func(args mock.Arguments) {
				f := args.Get(1).(func(context.Context) error)
				f(ctx)
			}).
			Return(errs.ErrInsufficientFundsError)

		uc := NewSendCoinUseCase(tm, tr, ur)
		err := uc.Send(ctx, &entities.TransferData{
			SenderID:  1,
			Recipient: "recipient",
			Amount:    100,
		})

		assert.ErrorIs(t, err, errs.ErrInsufficientFundsError)
		ur.AssertExpectations(t)
		tr.AssertExpectations(t)
		tm.AssertExpectations(t)
	})

	t.Run("get sender error", func(t *testing.T) {
		tm := new(MockTransactionManager)
		tr := new(MockTransactionRepo)
		ur := new(MockUserRepo)

		expectedErr := errors.New("sender error")
		ur.On("GetByID", ctx, 1).Return(nil, expectedErr)

		tm.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).
			Run(func(args mock.Arguments) {
				f := args.Get(1).(func(context.Context) error)
				f(ctx)
			}).
			Return(expectedErr)

		uc := NewSendCoinUseCase(tm, tr, ur)
		err := uc.Send(ctx, &entities.TransferData{
			SenderID:  1,
			Recipient: "recipient",
			Amount:    100,
		})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		ur.AssertExpectations(t)
	})

	t.Run("get recipient error", func(t *testing.T) {
		tm := new(MockTransactionManager)
		tr := new(MockTransactionRepo)
		ur := new(MockUserRepo)

		ur.On("GetByID", ctx, 1).Return(&entities.User{ID: 1}, nil)
		expectedErr := errors.New("recipient error")
		ur.On("GetByUsername", ctx, "recipient").Return(nil, expectedErr)

		tm.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).
			Run(func(args mock.Arguments) {
				f := args.Get(1).(func(context.Context) error)
				f(ctx)
			}).
			Return(expectedErr)

		uc := NewSendCoinUseCase(tm, tr, ur)
		err := uc.Send(ctx, &entities.TransferData{
			SenderID:  1,
			Recipient: "recipient",
			Amount:    100,
		})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		ur.AssertExpectations(t)
	})

	t.Run("get balance error", func(t *testing.T) {
		tm := new(MockTransactionManager)
		tr := new(MockTransactionRepo)
		ur := new(MockUserRepo)

		ur.On("GetByID", ctx, 1).Return(&entities.User{ID: 1}, nil)
		ur.On("GetByUsername", ctx, "recipient").Return(&entities.User{ID: 2}, nil)
		expectedErr := errors.New("balance error")
		tr.On("GetUserBalance", ctx, 1).Return(0, expectedErr)

		tm.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).
			Run(func(args mock.Arguments) {
				f := args.Get(1).(func(context.Context) error)
				f(ctx)
			}).
			Return(expectedErr)

		uc := NewSendCoinUseCase(tm, tr, ur)
		err := uc.Send(ctx, &entities.TransferData{
			SenderID:  1,
			Recipient: "recipient",
			Amount:    100,
		})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		ur.AssertExpectations(t)
		tr.AssertExpectations(t)
	})

	t.Run("create transaction error", func(t *testing.T) {
		tm := new(MockTransactionManager)
		tr := new(MockTransactionRepo)
		ur := new(MockUserRepo)

		ur.On("GetByID", ctx, 1).Return(&entities.User{ID: 1}, nil)
		ur.On("GetByUsername", ctx, "recipient").Return(&entities.User{ID: 2}, nil)
		tr.On("GetUserBalance", ctx, 1).Return(200, nil)

		expectedErr := errors.New("transaction error")
		tr.On("Create", ctx, mock.AnythingOfType("entities.TransactionData")).
			Return(nil, expectedErr)

		tm.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).
			Run(func(args mock.Arguments) {
				f := args.Get(1).(func(context.Context) error)
				f(ctx)
			}).
			Return(expectedErr)

		uc := NewSendCoinUseCase(tm, tr, ur)
		err := uc.Send(ctx, &entities.TransferData{
			SenderID:  1,
			Recipient: "recipient",
			Amount:    100,
		})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		ur.AssertExpectations(t)
		tr.AssertExpectations(t)
	})
}
func TestSend_CreateSecondTransactionError(t *testing.T) {
	ctx := t.Context()
	tm := new(MockTransactionManager)
	tr := new(MockTransactionRepo)
	ur := new(MockUserRepo)

	sender := &entities.User{ID: 1}
	recipient := &entities.User{ID: 2, Username: "recipient"}

	ur.On("GetByID", ctx, sender.ID).Return(sender, nil)
	ur.On("GetByUsername", ctx, recipient.Username).Return(recipient, nil)
	tr.On("GetUserBalance", ctx, sender.ID).Return(200, nil)

	tr.On("Create", ctx, mock.MatchedBy(func(data entities.TransactionData) bool {
		return data.UserID == sender.ID
	})).Return(&entities.Transaction{}, nil)

	expectedErr := errors.New("second transaction error")
	tr.On("Create", ctx, mock.MatchedBy(func(data entities.TransactionData) bool {
		return data.UserID == recipient.ID
	})).Return(nil, expectedErr)

	tm.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			f := args.Get(1).(func(context.Context) error)
			f(ctx)
		}).
		Return(expectedErr)

	uc := NewSendCoinUseCase(tm, tr, ur)
	err := uc.Send(ctx, &entities.TransferData{
		SenderID:  sender.ID,
		Recipient: recipient.Username,
		Amount:    100,
	})

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	ur.AssertExpectations(t)
	tr.AssertExpectations(t)
}
