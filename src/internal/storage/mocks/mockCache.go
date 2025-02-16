package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type MockCache struct {
	mock.Mock
}

func (m *MockCache) GetCachedRefreshToken(ctx context.Context, userID int) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockCache) CacheRefreshToken(ctx context.Context, userID int, token string) error {
	args := m.Called(ctx, userID, token)
	return args.Error(0)
}
