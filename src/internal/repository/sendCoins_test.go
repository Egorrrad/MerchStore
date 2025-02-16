package repository

import (
	"MerchStore/src/internal/storage/mocks"
	"MerchStore/src/internal/storage/model"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSendCoins_Success(t *testing.T) {
	mockStorage := new(mocks.MockStorage)
	mockTx := new(mocks.MockTransaction)
	repo := Repository{Storage: mockStorage}
	ctx := context.Background()

	sender := model.User{UserID: 1, Username: "user1", Coins: 100}
	receiver := model.User{UserID: 2, Username: "user2", Coins: 50}

	// Мокаем BeginTx.
	mockStorage.On("BeginTx", ctx).Return(mockTx, nil)

	// Мокаем GetUserForUpdate для отправителя и получателя.
	mockTx.On("GetUserForUpdate", ctx, "user1").Return(&sender, nil)
	mockTx.On("GetUserForUpdate", ctx, "user2").Return(&receiver, nil)

	// Мокаем UpdateUserCoins для отправителя и получателя.
	mockTx.On("UpdateUserCoins", ctx, sender.UserID, 50).Return(nil)    // 100 - 50 = 50
	mockTx.On("UpdateUserCoins", ctx, receiver.UserID, 100).Return(nil) // 50 + 50 = 100

	// Мокаем AddOperation.
	mockTx.On("AddOperation", ctx, sender.UserID, receiver.UserID, 50).Return(nil)

	// Мокаем Commit.
	mockTx.On("Commit").Return(nil)

	err := repo.SendCoins(ctx, "user1", "user2", 50)
	assert.NoError(t, err)

	// Проверяем, что все методы были вызваны.
	mockStorage.AssertExpectations(t)
	mockTx.AssertExpectations(t)
}

func TestSendCoins_NotEnoughCoins(t *testing.T) {
	mockStorage := new(mocks.MockStorage)
	mockTx := new(mocks.MockTransaction)
	repo := Repository{Storage: mockStorage}
	ctx := context.Background()

	sender := model.User{UserID: 1, Username: "user1", Coins: 30}
	receiver := model.User{UserID: 2, Username: "user2", Coins: 50}

	// Мокаем BeginTx.
	mockStorage.On("BeginTx", ctx).Return(mockTx, nil)

	// Мокаем GetUserForUpdate для отправителя и получателя.
	mockTx.On("GetUserForUpdate", ctx, "user1").Return(&sender, nil)
	mockTx.On("GetUserForUpdate", ctx, "user2").Return(&receiver, nil)

	// Мокаем Rollback.
	mockTx.On("Rollback").Return(nil)

	err := repo.SendCoins(ctx, "user1", "user2", 50)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not enough coins")
}

func TestSendCoins_ToYourself(t *testing.T) {
	mockStorage := new(mocks.MockStorage)
	repo := Repository{Storage: mockStorage}
	ctx := context.Background()

	// Мокаем BeginTx, чтобы убедиться, что он не вызывается.
	mockStorage.AssertNotCalled(t, "BeginTx")

	err := repo.SendCoins(ctx, "user1", "user1", 50)
	assert.Error(t, err)
	assert.Equal(t, ErrMsgSentToSelf, err)
}

func TestSendCoins_NegativeAmount(t *testing.T) {
	mockStorage := new(mocks.MockStorage)
	repo := Repository{Storage: mockStorage}
	ctx := context.Background()

	// Мокаем BeginTx, чтобы убедиться, что он не вызывается.
	mockStorage.AssertNotCalled(t, "BeginTx")

	// Пытаемся отправить отрицательное количество монет.
	err := repo.SendCoins(ctx, "user1", "user2", -50)
	assert.Error(t, err)
	assert.Equal(t, ErrMsgInvalidAmount, err)
}
