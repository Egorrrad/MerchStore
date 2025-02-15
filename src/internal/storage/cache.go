package storage

import (
	"MerchStore/src/internal/storage/cache"
	"context"
)

type CacheStorage interface {
	CacheRefreshToken(ctx context.Context, userID int, token string) error
	GetCachedRefreshToken(ctx context.Context, userID int) (string, error)
}

func NewCacheStorage(adr string) CacheStorage {
	storage := cache.InitRedis(adr)
	return storage
}
