package usecase

import (
	"av-merch-shop/internal/entities"
	"context"

	"github.com/stretchr/testify/mock"
)

type MockTransactionManager struct {
	mock.Mock
}

func (m *MockTransactionManager) Do(ctx context.Context, f func(ctx context.Context) error) error {
	args := m.Called(ctx, f)
	if f != nil {
		f(ctx)
	}
	return args.Error(0)
}

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) GetByID(ctx context.Context, userID int) (*entities.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepo) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepo) Create(ctx context.Context, data entities.UserData) (*entities.User, error) {
	args := m.Called(ctx, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

type MockTransactionRepo struct {
	mock.Mock
}

func (m *MockTransactionRepo) GetUserBalance(ctx context.Context, userID int) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

func (m *MockTransactionRepo) GetIncomingForUser(ctx context.Context, userID int) (*entities.UserReceived, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.UserReceived), args.Error(1)
}

func (m *MockTransactionRepo) GetOutgoingForUser(ctx context.Context, userID int) (*entities.UserSent, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.UserSent), args.Error(1)
}

func (m *MockTransactionRepo) Create(ctx context.Context, data entities.TransactionData) (*entities.Transaction, error) {
	args := m.Called(ctx, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Transaction), args.Error(1)
}

type MockItemRepo struct {
	mock.Mock
}

func (m *MockItemRepo) GetItemByName(ctx context.Context, itemName string) (*entities.Item, error) {
	args := m.Called(ctx, itemName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Item), args.Error(1)
}

type MockOrderRepo struct {
	mock.Mock
}

func (m *MockOrderRepo) Create(ctx context.Context, data entities.OrderData) (*entities.Order, error) {
	args := m.Called(ctx, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Order), args.Error(1)
}

func (m *MockOrderRepo) GetUserInventory(ctx context.Context, userID int) (*entities.UserInventory, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.UserInventory), args.Error(1)
}

type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) GenerateToken(userID int, username, role string) (string, error) {
	args := m.Called(userID, username, role)
	return args.String(0), args.Error(1)
}

type MockHashService struct {
	mock.Mock
}

func (m *MockHashService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockHashService) CompareWithPassword(hashed, password string) bool {
	args := m.Called(hashed, password)
	return args.Bool(0)
}

type MockUserInfoCache struct {
	mock.Mock
}

func (m *MockUserInfoCache) GetUserInfo(ctx context.Context, userID int) (*entities.UserInfo, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.UserInfo), args.Error(1)
}

func (m *MockUserInfoCache) SetUserInfo(ctx context.Context, userID int, info entities.UserInfo) error {
	args := m.Called(ctx, userID, info)
	return args.Error(0)
}

func (m *MockUserInfoCache) ExpireUserInfo(ctx context.Context, userID int) {
	m.Called(ctx, userID)
}
