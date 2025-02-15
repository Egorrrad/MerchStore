package postgres

import (
	"MerchStore/src/internal/storage/model"
	"context"
)

func (p *Storage) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	query := `SELECT user_id, username, password_hash, coins FROM users WHERE username = $1`
	err := p.DB.QueryRowContext(ctx, query, username).Scan(&user.UserID, &user.Username, &user.PasswordHash, &user.Coins)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
