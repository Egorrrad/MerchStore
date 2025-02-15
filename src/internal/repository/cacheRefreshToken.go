package repository

import (
	"context"
	"time"
)

func (r Repository) SaveRefreshToken(ctx context.Context, userID int, token string, expiresAt time.Time) error {
	// Сохраняем в PostgreSQL
	if err := r.storage.SaveRefreshToken(ctx, userID, token, expiresAt); err != nil {
		return err
	}

	// Кэшируем в Redis
	return r.cache.CacheRefreshToken(ctx, userID, token)
}
