package repository

import (
	"context"
)

func (r Repository) BuyItem(ctx context.Context, username, productName string) error {
	user, err := r.storage.GetUserByUsername(ctx, username)
	if err != nil {
		return err
	}

	product, err := r.storage.GetProduct(ctx, productName)
	if err != nil {
		return err
	}
	if product.Quantity == 0 {
		return ErrMsgOutOfStock
	}
	if user.Coins < product.Price {
		return ErrMsgNotEnoughCoins
	}

	err = r.storage.UpdateUserCoins(ctx, user.UserID, user.Coins-product.Price)
	if err != nil {
		return err
	}
	err = r.storage.UpdateProductQuantity(ctx, product.ProductID, product.Quantity-1)
	if err != nil {
		return err
	}
	err = r.storage.AddPurchase(ctx, user.UserID, product.ProductID, 1)
	if err != nil {
		return err
	}
	return nil
}
