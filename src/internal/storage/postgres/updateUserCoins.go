package postgres

import (
	"context"
	"fmt"
)

func (p *Storage) UpdateUserCoins(ctx context.Context, userID int, coins int) error {
	query := `UPDATE users SET coins = $1 WHERE user_id = $2`

	_, err := p.DB.ExecContext(ctx, query, coins, userID)
	if err != nil {
		return fmt.Errorf("failed to update user coins: %w", err)
	}

	return nil
}
