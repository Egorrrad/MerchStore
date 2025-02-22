package repository

import (
	"MerchStore/src/internal/storage/mocks"
	"MerchStore/src/internal/storage/model"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBuyItem_Success(t *testing.T) {
	mockStorage := new(mocks.MockStorage)
	mockTx := new(mocks.MockTransaction)
	repo := Repository{Storage: mockStorage}
	ctx := context.Background()

	user := model.User{UserID: 1, Username: "testuser", Coins: 100}
	product := model.Product{ProductID: 1, Name: "item1", Price: 50, Quantity: 10}

	// Мокаем BeginTx.
	mockStorage.On("BeginTx", ctx).Return(mockTx, nil)

	// Мокаем методы транзакции.
	mockTx.On("GetUserForUpdate", ctx, "testuser").Return(&user, nil)
	mockTx.On("GetProductForUpdate", ctx, "item1").Return(&product, nil)
	mockTx.On("UpdateUserCoins", ctx, user.UserID, 50).Return(nil)
	mockTx.On("UpdateProductQuantity", ctx, product.ProductID, 9).Return(nil)
	mockTx.On("AddPurchase", ctx, user.UserID, product.ProductID, 1).Return(nil)
	mockTx.On("Commit").Return(nil) // Ожидаем Commit.
	// Не добавляем Rollback, так как он не должен вызываться в успешном сценарии.

	err := repo.BuyItem(ctx, "testuser", "item1")
	assert.NoError(t, err)

	// Проверяем, что все методы были вызваны.
	mockStorage.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestBuyItem_UserNotFound(t *testing.T) {
	mockStorage := new(mocks.MockStorage)
	mockTx := new(mocks.MockTransaction)
	repo := Repository{Storage: mockStorage}
	ctx := context.Background()

	// Мокаем BeginTx.
	mockStorage.On("BeginTx", ctx).Return(mockTx, nil)

	// Мокаем GetUserForUpdate с ошибкой.
	mockTx.On("GetUserForUpdate", ctx, "testuser").Return(&model.User{}, errors.New("user not found"))
	mockTx.On("Rollback").Return(nil) // Ожидаем Rollback, так как будет ошибка.

	err := repo.BuyItem(ctx, "testuser", "item1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "get user error") // Проверяем, что ошибка содержит текст "get user error".

	// Проверяем, что все методы были вызваны.
	mockStorage.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestBuyItem_ProductNotFound(t *testing.T) {
	mockStorage := new(mocks.MockStorage)
	mockTx := new(mocks.MockTransaction)
	repo := Repository{Storage: mockStorage}
	ctx := context.Background()

	user := model.User{UserID: 1, Username: "testuser", Coins: 100}

	// Мокаем BeginTx.
	mockStorage.On("BeginTx", ctx).Return(mockTx, nil)

	// Мокаем GetUserForUpdate.
	mockTx.On("GetUserForUpdate", ctx, "testuser").Return(&user, nil)

	// Мокаем GetProductForUpdate с ошибкой.
	mockTx.On("GetProductForUpdate", ctx, "item1").Return(&model.Product{}, errors.New("product not found"))
	mockTx.On("Rollback").Return(nil) // Ожидаем Rollback, так как будет ошибка.

	err := repo.BuyItem(ctx, "testuser", "item1")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "product not exist") // Проверяем, что ошибка содержит текст "get product error".

	// Проверяем, что все методы были вызваны.
	mockStorage.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}
func TestBuyItem_OutOfStock(t *testing.T) {
	mockStorage := new(mocks.MockStorage)
	mockTx := new(mocks.MockTransaction)
	repo := Repository{Storage: mockStorage}
	ctx := context.Background()

	user := model.User{UserID: 1, Username: "testuser", Coins: 100}
	product := model.Product{ProductID: 1, Name: "item1", Price: 50, Quantity: 0} // Товара нет в наличии.

	// Мокаем BeginTx.
	mockStorage.On("BeginTx", ctx).Return(mockTx, nil)

	// Мокаем GetUserForUpdate и GetProductForUpdate.
	mockTx.On("GetUserForUpdate", ctx, "testuser").Return(&user, nil)
	mockTx.On("GetProductForUpdate", ctx, "item1").Return(&product, nil)
	mockTx.On("Rollback").Return(nil) // Ожидаем Rollback, так как будет ошибка.

	err := repo.BuyItem(ctx, "testuser", "item1")
	assert.Error(t, err)
	assert.Equal(t, ErrMsgOutOfStock, err) // Проверяем, что возвращается ошибка "out of stock".

	// Проверяем, что все методы были вызваны.
	mockStorage.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestBuyItem_NotEnoughCoins(t *testing.T) {
	mockStorage := new(mocks.MockStorage)
	mockTx := new(mocks.MockTransaction)
	repo := Repository{Storage: mockStorage}
	ctx := context.Background()

	user := model.User{UserID: 1, Username: "testuser", Coins: 30} // Недостаточно монет.
	product := model.Product{ProductID: 1, Name: "item1", Price: 50, Quantity: 10}

	// Мокаем BeginTx.
	mockStorage.On("BeginTx", ctx).Return(mockTx, nil)

	// Мокаем GetUserForUpdate и GetProductForUpdate.
	mockTx.On("GetUserForUpdate", ctx, "testuser").Return(&user, nil)
	mockTx.On("GetProductForUpdate", ctx, "item1").Return(&product, nil)
	mockTx.On("Rollback").Return(nil) // Ожидаем Rollback, так как будет ошибка.

	err := repo.BuyItem(ctx, "testuser", "item1")
	assert.Error(t, err)
	assert.Equal(t, ErrMsgNotEnoughCoins, err) // Проверяем, что возвращается ошибка "not enough coins".

	// Проверяем, что все методы были вызваны.
	mockStorage.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}
