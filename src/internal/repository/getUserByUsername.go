package repository

import (
	"MerchStore/src/internal/datastorage/model"
	"context"
	"database/sql"
)

func (r Repository) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	user, err := r.storage.GetUserByUsername(ctx, username)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}
