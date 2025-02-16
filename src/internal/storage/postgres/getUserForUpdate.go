package postgres

import (
	"MerchStore/src/internal/storage/model"
	"context"
	"fmt"
)

func (t *StorageTx) GetUserForUpdate(ctx context.Context, username string) (*model.User, error) {
	query := `
        SELECT user_id, username, coins 
        FROM users 
        WHERE username = $1 
        FOR UPDATE` // Блокировка строки

	row := t.tx.QueryRowContext(ctx, query, username)
	var user model.User
	if err := row.Scan(&user.UserID, &user.Username, &user.Coins); err != nil {
		return nil, fmt.Errorf("get user for update: %w", err)
	}
	return &user, nil
}
