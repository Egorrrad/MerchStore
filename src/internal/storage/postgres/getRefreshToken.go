package postgres

import (
	"MerchStore/src/internal/storage/model"
	"context"
)

func (p *Storage) GetRefreshToken(ctx context.Context, userID int, token string) (*model.RefreshToken, error) {
	var rt model.RefreshToken
	query := `SELECT token_id, user_id, token, expires_at, created_at FROM refresh_tokens WHERE user_id = $1 AND token = $2`
	err := p.DB.QueryRowContext(ctx, query, userID, token).Scan(&rt.TokenID, &rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.CreatedAt)
	return &rt, err
}
