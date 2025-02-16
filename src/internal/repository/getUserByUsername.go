package repository

import (
	"MerchStore/src/internal/storage/model"
	"context"
	"database/sql"
)

func (r Repository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	user, err := r.Storage.GetUserByUsername(ctx, username)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}
