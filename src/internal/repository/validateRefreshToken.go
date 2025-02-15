package repository

import (
	"context"
)

func (r Repository) ValidateRefreshToken(ctx context.Context, userID int, token string) error {
	// Сначала проверяем Redis
	cachedToken, err := r.cache.GetCachedRefreshToken(ctx, userID)
	if err == nil && cachedToken == token {
		return nil
	}

	// Если нет в Redis, проверяем PostgreSQL
	_, err = r.storage.GetRefreshToken(ctx, userID, token)
	if err != nil {
		return err
	}

	// Обновляем кэш при успехе
	err = r.cache.CacheRefreshToken(ctx, userID, token)
	if err != nil {
		return err
	}

	return nil
}
