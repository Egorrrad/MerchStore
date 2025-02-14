package repository

import "context"

func (r Repository) SendCoins(ctx context.Context, fromUser, toUser string, amount int) error {
	sender, err := r.storage.GetUserByUsername(ctx, fromUser)
	if err != nil {
		return err
	}
	resiever, err := r.storage.GetUserByUsername(ctx, toUser)
	if err != nil {
		return err
	}
	if sender.Coins < amount {
		return ErrMsgNotEnoughCoins
	}
	err = r.storage.UpdateUserCoins(ctx, sender.UserID, sender.Coins-amount)
	if err != nil {
		return err
	}
	err = r.storage.UpdateUserCoins(ctx, resiever.UserID, resiever.Coins+amount)
	if err != nil {
		return err
	}
	err = r.storage.AddOperation(ctx, sender.UserID, resiever.UserID, amount)
	if err != nil {
		return err
	}
	return nil
}
