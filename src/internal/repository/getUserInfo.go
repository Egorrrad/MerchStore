package repository

import (
	"MerchStore/src/internal/storage/model"
	"context"
)

func (r Repository) GetUserInfo(ctx context.Context, username string) (
	*model.User,
	[]model.Purchase,
	[]model.Operation,
	error) {
	user, err := r.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, nil, nil, ErrMsgUserNotExist
	}

	purchases, err := r.Storage.GetUserPurchases(ctx, user.UserID)
	if err != nil {
		return nil, nil, nil, err
	}

	operations, err := r.Storage.GetUserOperations(ctx, user.UserID)
	if err != nil {
		return nil, nil, nil, err
	}

	return user, purchases, operations, nil
}
