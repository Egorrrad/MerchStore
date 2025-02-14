package repository

import (
	"MerchStore/src/internal/datastorage/model"
	"context"
)

func (r *Repository) GetUserInfo(ctx context.Context, username string) (
	[]model.Purchase,
	[]model.Operation,
	error) {
	user, err := r.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, nil, ErrMsgUserNotExist
	}

	purchases, err := r.storage.GetUserPurchases(ctx, user.UserID)
	if err != nil {
		return nil, nil, err
	}

	operations, err := r.storage.GetUserOperations(ctx, user.UserID)
	if err != nil {
		return nil, nil, err
	}

	return purchases, operations, nil
}
