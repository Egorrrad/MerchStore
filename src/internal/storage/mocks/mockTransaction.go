package mocks

import (
	"MerchStore/src/internal/storage/model"
	"context"
)

func (m *MockTransaction) GetUserForUpdate(ctx context.Context, username string) (*model.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockTransaction) GetProductForUpdate(ctx context.Context, productName string) (*model.Product, error) {
	args := m.Called(ctx, productName)
	return args.Get(0).(*model.Product), args.Error(1)
}

func (m *MockTransaction) UpdateUserCoins(ctx context.Context, userID int, newAmount int) error {
	args := m.Called(ctx, userID, newAmount)
	return args.Error(0)
}

func (m *MockTransaction) UpdateProductQuantity(ctx context.Context, productID int, quantity int) error {
	args := m.Called(ctx, productID, quantity)
	return args.Error(0)
}

func (m *MockTransaction) AddOperation(ctx context.Context, from, to int, amount int) error {
	args := m.Called(ctx, from, to, amount)
	return args.Error(0)
}

func (m *MockTransaction) AddPurchase(ctx context.Context, userID, productID, quantity int) error {
	args := m.Called(ctx, userID, productID, quantity)
	return args.Error(0)
}

func (m *MockTransaction) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTransaction) Rollback() error {
	args := m.Called()
	return args.Error(0)
}
