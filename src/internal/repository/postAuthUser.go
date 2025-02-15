package repository

import (
	"MerchStore/src/internal/auth"
	"context"
	"database/sql"
	"time"
)

func (r Repository) PostAuthUser(ctx context.Context, username, password string) (*string, error) {
	user, err := r.storage.GetUserByUsername(ctx, username)
	if err != nil && err == sql.ErrNoRows {
		err1 := r.storage.AddUser(ctx, username, password, "user")
		if err1 != nil {
			return nil, err1
		}
	} else {
		return nil, err
	}
	var userID int
	if user == nil {
		user, err = r.storage.GetUserByUsername(ctx, username)
		if err != nil {
			return nil, err
		}
	}
	userID = user.UserID
	// Генерация токена
	token, err := auth.GenerateRefreshToken(userID, user.Username)
	if err != nil {
		return nil, ErrMsgTokenGenFailed
	}

	// Сохранение токена в репозитории
	if err = r.SaveRefreshToken(ctx, userID, token, time.Now().Add(7*24*time.Hour)); err != nil {
		return nil, ErrMsgTokenSaveFailed
	}

	return &token, nil
}
