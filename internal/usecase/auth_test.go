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

func TestSignIn_Errors(t *testing.T) {
	ctx := context.Background()

	t.Run("transaction manager error", func(t *testing.T) {
		tm := new(MockTransactionManager)
		ur := new(MockUserRepo)
		tr := new(MockTransactionRepo)
		ts := new(MockTokenService)
		hs := new(MockHashService)

		ur.On("GetByUsername", ctx, "user").Return(nil, errors.New("some error"))

		tm.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).Return(errors.New("transaction error"))

		uc := NewAuthUsecase(tm, ur, tr, ts, hs)
		token, err := uc.SignIn(ctx, "user", "pass")

		assert.Error(t, err)
		assert.Empty(t, token)
		tm.AssertExpectations(t)
		ur.AssertExpectations(t)
	})

	t.Run("user repo error", func(t *testing.T) {
		tm := new(MockTransactionManager)
		ur := new(MockUserRepo)
		tr := new(MockTransactionRepo)
		ts := new(MockTokenService)
		hs := new(MockHashService)

		ur.On("GetByUsername", ctx, "user").Return(nil, errors.New("db error"))

		tm.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).Return(errors.New("db error"))

		uc := NewAuthUsecase(tm, ur, tr, ts, hs)
		token, err := uc.SignIn(ctx, "user", "pass")

		assert.Error(t, err)
		assert.Empty(t, token)
		tm.AssertExpectations(t)
		ur.AssertExpectations(t)
	})

	t.Run("hash service error", func(t *testing.T) {
		tm := new(MockTransactionManager)
		ur := new(MockUserRepo)
		tr := new(MockTransactionRepo)
		ts := new(MockTokenService)
		hs := new(MockHashService)

		ur.On("GetByUsername", ctx, "user").Return(nil, errs.ErrNotFoundError)

		hs.On("HashPassword", "pass").Return("", errors.New("hash error"))

		tm.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).Return(errors.New("hash error"))

		uc := NewAuthUsecase(tm, ur, tr, ts, hs)
		token, err := uc.SignIn(ctx, "user", "pass")

		assert.Error(t, err)
		assert.Empty(t, token)
		tm.AssertExpectations(t)
		ur.AssertExpectations(t)
		hs.AssertExpectations(t)
	})

	t.Run("user creation error", func(t *testing.T) {
		tm := new(MockTransactionManager)
		ur := new(MockUserRepo)
		tr := new(MockTransactionRepo)
		ts := new(MockTokenService)
		hs := new(MockHashService)

		ur.On("GetByUsername", ctx, "user").Return(nil, errs.ErrNotFoundError)

		hs.On("HashPassword", "pass").Return("hashed", nil)

		ur.On("Create", ctx, mock.AnythingOfType("entities.UserData")).Return(nil, errors.New("creation error"))

		tm.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).Return(errors.New("creation error"))

		uc := NewAuthUsecase(tm, ur, tr, ts, hs)
		token, err := uc.SignIn(ctx, "user", "pass")

		assert.Error(t, err)
		assert.Empty(t, token)
		tm.AssertExpectations(t)
		ur.AssertExpectations(t)
		hs.AssertExpectations(t)
	})
}

func TestGetTokenForNewUser_TokenServiceError(t *testing.T) {
	ctx := context.Background()
	tm := new(MockTransactionManager)
	ur := new(MockUserRepo)
	tr := new(MockTransactionRepo)
	ts := new(MockTokenService)
	hs := new(MockHashService)

	ur.On("GetByUsername", ctx, "user").Return(nil, errs.ErrNotFoundError)
	hs.On("HashPassword", "pass").Return("hashed", nil)

	newUser := &entities.User{
		ID:       1,
		Username: "user",
		Role:     entities.RoleUser,
	}
	ur.On("Create", ctx, mock.AnythingOfType("entities.UserData")).Return(newUser, nil)
	tr.On("Create", ctx, mock.AnythingOfType("entities.TransactionData")).Return(&entities.Transaction{}, nil)

	expectedErr := errors.New("token generation error")
	ts.On("GenerateToken", newUser.ID, newUser.Username, newUser.Role).Return("", expectedErr)

	tm.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			f := args.Get(1).(func(context.Context) error)
			f(ctx)
		}).
		Return(expectedErr)

	uc := NewAuthUsecase(tm, ur, tr, ts, hs)
	token, err := uc.SignIn(ctx, "user", "pass")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, expectedErr, err)
	ur.AssertExpectations(t)
	hs.AssertExpectations(t)
	tr.AssertExpectations(t)
	ts.AssertExpectations(t)
}

func TestSignIn_CreditInitialAmountError(t *testing.T) {
	ctx := context.Background()
	tm := new(MockTransactionManager)
	ur := new(MockUserRepo)
	tr := new(MockTransactionRepo)
	ts := new(MockTokenService)
	hs := new(MockHashService)

	ur.On("GetByUsername", ctx, "user").Return(nil, errs.ErrNotFoundError)
	hs.On("HashPassword", "pass").Return("hashed", nil)

	newUser := &entities.User{
		ID:       1,
		Username: "user",
		Role:     entities.RoleUser,
	}
	ur.On("Create", ctx, mock.AnythingOfType("entities.UserData")).Return(newUser, nil)

	expectedErr := errors.New("credit transaction error")
	tr.On("Create", ctx, mock.MatchedBy(func(data entities.TransactionData) bool {
		return data.UserID == newUser.ID &&
			data.Amount == DefaultBalanceForUser &&
			data.TransactionType == entities.TransactionCredit
	})).Return(nil, expectedErr)

	tm.On("Do", ctx, mock.AnythingOfType("func(context.Context) error")).
		Run(func(args mock.Arguments) {
			f := args.Get(1).(func(context.Context) error)
			f(ctx)
		}).
		Return(expectedErr)

	uc := NewAuthUsecase(tm, ur, tr, ts, hs)
	token, err := uc.SignIn(ctx, "user", "pass")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, expectedErr, err)
	ur.AssertExpectations(t)
	hs.AssertExpectations(t)
	tr.AssertExpectations(t)
}
