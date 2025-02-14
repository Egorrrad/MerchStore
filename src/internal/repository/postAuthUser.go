package repository

import (
	"MerchStore/src/internal/auth"
	"context"
	"database/sql"
)

func (r *Repository) PostAuthUser(ctx context.Context, username, password string) (*string, error) {
	_, err := r.storage.GetUserByUsername(ctx, username)
	if err == sql.ErrNoRows {
		err1 := r.storage.AddUser(ctx, username, password, "user")
		if err1 != nil {
			return nil, err1
		}
	} else {
		return nil, err
	}

	// Генерируем JWT токен
	token, err := auth.GenerateJWT(username)
	if err != nil {
		return nil, err
	}
	return &token, nil
}
