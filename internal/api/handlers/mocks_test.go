package handlers

import (
	"av-merch-shop/internal/entities"
	"context"

	"github.com/stretchr/testify/mock"
)

type mockAuthUseCase struct {
	mock.Mock
}

func (m *mockAuthUseCase) SignIn(ctx context.Context, username string, password string) (string, error) {
	args := m.Called(ctx, username, password)
	return args.String(0), args.Error(1)
}

type mockInfoUseCase struct {
	mock.Mock
}

func (m *mockInfoUseCase) GetInfo(ctx context.Context, userID int) (*entities.UserInfo, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.UserInfo), args.Error(1)
}

type mockSendCoinUseCase struct {
	mock.Mock
}

func (m *mockSendCoinUseCase) Send(ctx context.Context, data *entities.TransferData) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

type mockOrderUseCase struct {
	mock.Mock
}

func (m *mockOrderUseCase) Buy(ctx context.Context, data *entities.OrderRequest) error {
	args := m.Called(ctx, data)
	return args.Error(0)
}

type mockCreditUseCase struct {
	mock.Mock
}

func (m *mockCreditUseCase) Credit(ctx context.Context, data *entities.CreditData) (*entities.CreditTransactionResult, error) {
	args := m.Called(ctx, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.CreditTransactionResult), args.Error(1)
}

type mockValidator struct {
	mock.Mock
}

func (m *mockValidator) Engine() interface{} {
	args := m.Called()
	return args.Get(0)
}

func (m *mockValidator) ValidateStruct(interface{}) error {
	args := m.Called()
	return args.Error(0)
}
