package repository

import (
	"MerchStore/src/internal/auth"
	"context"
	"time"
)

func (r Repository) ValidateRefreshToken(ctx context.Context, token, secretKey string) (*string, error) {
	// Парсинг токена
	claims, err := auth.ParseRefreshToken(token, secretKey)
	if err != nil {
		return nil, ErrMsgInvalidToken
	}

	// Извлекаем user_id, username и exp из токена
	userID := int(claims["user_id"].(float64))
	username := claims["username"].(string)
	exp := claims["exp"].(float64)

	// Проверяем истечение срока действия
	if time.Now().Unix() > int64(exp) {
		return nil, ErrMsgTokenExpired
	}
	// Сначала проверяем Redis
	cachedToken, err := r.Cache.GetCachedRefreshToken(ctx, userID)
	if err == nil && cachedToken == token {
		return &username, nil
	}

	// Если нет в Redis, проверяем PostgreSQL
	_, err = r.Storage.GetRefreshToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Обновляем кэш при успехе
	err = r.Cache.CacheRefreshToken(ctx, userID, token)
	if err != nil {
		return nil, err
	}

	return &username, nil
}
