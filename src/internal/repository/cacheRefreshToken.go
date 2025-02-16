package repository

import (
	"context"
	"time"
)

func (r Repository) SaveRefreshToken(ctx context.Context, userID int, token string, expiresAt time.Time) error {
	// Сохраняем в PostgreSQL
	if err := r.Storage.SaveRefreshToken(ctx, userID, token, expiresAt); err != nil {
		return err
	}

	// Кэшируем в Redis
	return r.Cache.CacheRefreshToken(ctx, userID, token)
}
