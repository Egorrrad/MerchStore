package postgres

import (
	"context"
	"fmt"
)

func (t *StorageTx) UpdateUserCoins(ctx context.Context, userID int, newAmount int) error {
	query := `UPDATE users SET coins = $1 WHERE user_id = $2`
	_, err := t.tx.ExecContext(ctx, query, newAmount, userID)
	if err != nil {
		return fmt.Errorf("failed to update user coins: %w", err)
	}
	return nil
}
