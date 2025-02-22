package postgres

import (
	"MerchStore/src/internal/storage/model"
	"context"
)

func (p *Storage) GetRefreshToken(ctx context.Context, userID int) (*model.RefreshToken, error) {
	var rt model.RefreshToken
	query := `SELECT token_id, user_id, token, expires_at, created_at FROM refresh_tokens WHERE user_id = $1`
	err := p.DB.QueryRowContext(ctx, query, userID).Scan(&rt.TokenID, &rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &rt, err
}
