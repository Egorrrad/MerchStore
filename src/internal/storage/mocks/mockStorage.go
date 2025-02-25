package mocks

import (
	"MerchStore/src/internal/storage/model"
	"MerchStore/src/internal/storage/postgres"
	"context"
	"github.com/stretchr/testify/mock"
	"time"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) UpdateRefreshToken(ctx context.Context, id int, token string, expires time.Time) error {
	args := m.Called(ctx, id, token, expires)
	return args.Error(0)
}

type MockTransaction struct {
	mock.Mock
}

func (m *MockStorage) BeginTx(ctx context.Context) (postgres.Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(postgres.Tx), args.Error(1)
}

func (m *MockStorage) AddUser(ctx context.Context, username, passwordHash, role string) error {
	args := m.Called(ctx, username, passwordHash, role)
	return args.Error(0)
}

func (m *MockStorage) AddOperation(ctx context.Context, senderID, receiverID, amount int) error {
	args := m.Called(ctx, senderID, receiverID, amount)
	return args.Error(0)
}

func (m *MockStorage) GetUserOperations(ctx context.Context, userID int) ([]model.Operation, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.Operation), args.Error(1)
}

func (m *MockStorage) GetUserPurchases(ctx context.Context, userID int) ([]model.Purchase, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.Purchase), args.Error(1)
}

func (m *MockStorage) SaveRefreshToken(ctx context.Context, userID int, token string, expiresAt time.Time) error {
	args := m.Called(ctx, userID, token, expiresAt)
	return args.Error(0)
}

func (m *MockStorage) GetRefreshToken(ctx context.Context, userID int) (*model.RefreshToken, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*model.RefreshToken), args.Error(1)
}

func (m *MockStorage) DeleteRefreshToken(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockStorage) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockStorage) GetProduct(ctx context.Context, productName string) (*model.Product, error) {
	args := m.Called(ctx, productName)
	return args.Get(0).(*model.Product), args.Error(1)
}

func (m *MockStorage) UpdateUserCoins(ctx context.Context, userID int, coins int) error {
	args := m.Called(ctx, userID, coins)
	return args.Error(0)
}

func (m *MockStorage) UpdateProductQuantity(ctx context.Context, productID int, quantity int) error {
	args := m.Called(ctx, productID, quantity)
	return args.Error(0)
}

func (m *MockStorage) AddPurchase(ctx context.Context, userID, productID, quantity int) error {
	args := m.Called(ctx, userID, productID, quantity)
	return args.Error(0)
}
