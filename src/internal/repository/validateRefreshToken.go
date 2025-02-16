package repository

import (
	"context"
)

func (r Repository) ValidateRefreshToken(ctx context.Context, userID int, token string) error {
	// Сначала проверяем Redis
	cachedToken, err := r.Cache.GetCachedRefreshToken(ctx, userID)
	if err == nil && cachedToken == token {
		return nil
	}

	// Если нет в Redis, проверяем PostgreSQL
	_, err = r.Storage.GetRefreshToken(ctx, userID, token)
	if err != nil {
		return err
	}

	// Обновляем кэш при успехе
	err = r.Cache.CacheRefreshToken(ctx, userID, token)
	if err != nil {
		return err
	}

	return nil
}
