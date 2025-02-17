package repository

import (
	"MerchStore/src/internal/storage/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const testToken = "test_token"

func TestSaveRefreshToken_Success(t *testing.T) {
	mockStorage := new(mocks.MockStorage)
	mockCache := new(mocks.MockCache)
	repo := Repository{Storage: mockStorage, Cache: mockCache}
	ctx := context.Background()

	token := testToken
	expiresAt := time.Now().Add(24 * time.Hour)
	userID := 1

	mockStorage.On("SaveRefreshToken", ctx, userID, token, expiresAt).Return(nil)
	mockCache.On("CacheRefreshToken", ctx, userID, token).Return(nil)

	err := repo.SaveRefreshToken(ctx, userID, token, expiresAt)
	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestSaveRefreshToken_StorageError(t *testing.T) {
	mockStorage := new(mocks.MockStorage)
	mockCache := new(mocks.MockCache)
	repo := Repository{Storage: mockStorage, Cache: mockCache}
	ctx := context.Background()

	token := testToken
	expiresAt := time.Now().Add(24 * time.Hour)
	userID := 1

	mockStorage.On("SaveRefreshToken", ctx, userID, token, expiresAt).Return(errors.New("storage error"))

	err := repo.SaveRefreshToken(ctx, userID, token, expiresAt)
	assert.Error(t, err)
	assert.Equal(t, "storage error", err.Error())
	mockStorage.AssertExpectations(t)
}

func TestSaveRefreshToken_CacheError(t *testing.T) {
	mockStorage := new(mocks.MockStorage)
	mockCache := new(mocks.MockCache)
	repo := Repository{Storage: mockStorage, Cache: mockCache}
	ctx := context.Background()

	token := testToken
	expiresAt := time.Now().Add(24 * time.Hour)
	userID := 1

	mockStorage.On("SaveRefreshToken", ctx, userID, token, expiresAt).Return(nil)
	mockCache.On("CacheRefreshToken", ctx, userID, token).Return(errors.New("cache error"))

	err := repo.SaveRefreshToken(ctx, userID, token, expiresAt)
	assert.Error(t, err)
	assert.Equal(t, "cache error", err.Error())
	mockStorage.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}
