package tests

import (
	"MerchStore/src/internal/repository"
	"MerchStore/src/internal/storage/mocks"
	"MerchStore/src/internal/storage/model"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuyItem_Success(t *testing.T) {
	mockStorage := new(mocks.MockStorage)
	repo := repository.Repository{Storage: mockStorage}
	ctx := context.Background()

	user := model.User{UserID: 1, Username: "testuser", Coins: 100}
	product := model.Product{ProductID: 1, Name: "item1", Price: 50, Quantity: 10}

	mockStorage.On("GetUserByUsername", ctx, "testuser").Return(&user, nil)
	mockStorage.On("GetProduct", ctx, "item1").Return(&product, nil)
	mockStorage.On("UpdateUserCoins", ctx, user.UserID, 50).Return(nil)
	mockStorage.On("UpdateProductQuantity", ctx, product.ProductID, 9).Return(nil)
	mockStorage.On("AddPurchase", ctx, user.UserID, product.ProductID, 1).Return(nil)

	err := repo.BuyItem(ctx, "testuser", "item1")
	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
}

func TestBuyItem_UserNotFound(t *testing.T) {
	mockStorage := new(mocks.MockStorage)
	repo := repository.Repository{Storage: mockStorage}
	ctx := context.Background()

	mockStorage.On("GetUserByUsername", ctx, "testuser").Return(&model.User{}, errors.New("user not found"))

	err := repo.BuyItem(ctx, "testuser", "item1")
	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
}

func TestBuyItem_ProductNotFound(t *testing.T) {
	mockStorage := new(mocks.MockStorage)
	repo := repository.Repository{Storage: mockStorage}
	ctx := context.Background()

	user := model.User{UserID: 1, Username: "testuser", Coins: 100}
	mockStorage.On("GetUserByUsername", ctx, "testuser").Return(&user, nil)
	mockStorage.On("GetProduct", ctx, "item1").Return(&model.Product{}, errors.New("product not found"))

	err := repo.BuyItem(ctx, "testuser", "item1")
	assert.Error(t, err)
	assert.Equal(t, "product not found", err.Error())
}

func TestBuyItem_OutOfStock(t *testing.T) {
	mockStorage := new(mocks.MockStorage)
	repo := repository.Repository{Storage: mockStorage}
	ctx := context.Background()

	user := model.User{UserID: 1, Username: "testuser", Coins: 100}
	product := model.Product{ProductID: 1, Name: "item1", Price: 50, Quantity: 0}

	mockStorage.On("GetUserByUsername", ctx, "testuser").Return(&user, nil)
	mockStorage.On("GetProduct", ctx, "item1").Return(&product, nil)

	err := repo.BuyItem(ctx, "testuser", "item1")
	assert.Error(t, err)
	assert.Equal(t, repository.ErrMsgOutOfStock, err)
}

func TestBuyItem_NotEnoughCoins(t *testing.T) {
	mockStorage := new(mocks.MockStorage)
	repo := repository.Repository{Storage: mockStorage}
	ctx := context.Background()

	user := model.User{UserID: 1, Username: "testuser", Coins: 30}
	product := model.Product{ProductID: 1, Name: "item1", Price: 50, Quantity: 10}

	mockStorage.On("GetUserByUsername", ctx, "testuser").Return(&user, nil)
	mockStorage.On("GetProduct", ctx, "item1").Return(&product, nil)

	err := repo.BuyItem(ctx, "testuser", "item1")
	assert.Error(t, err)
	assert.Equal(t, repository.ErrMsgNotEnoughCoins, err)
}
