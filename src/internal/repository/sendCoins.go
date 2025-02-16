package repository

import "context"

func (r Repository) SendCoins(ctx context.Context, fromUser, toUser string, amount int) error {
	sender, err := r.Storage.GetUserByUsername(ctx, fromUser)
	if err != nil {
		return err
	}
	resiever, err := r.Storage.GetUserByUsername(ctx, toUser)
	if err != nil {
		return err
	}
	if sender.Coins < amount {
		return ErrMsgNotEnoughCoins
	}
	err = r.Storage.UpdateUserCoins(ctx, sender.UserID, sender.Coins-amount)
	if err != nil {
		return err
	}
	err = r.Storage.UpdateUserCoins(ctx, resiever.UserID, resiever.Coins+amount)
	if err != nil {
		return err
	}
	err = r.Storage.AddOperation(ctx, sender.UserID, resiever.UserID, amount)
	if err != nil {
		return err
	}
	return nil
}
