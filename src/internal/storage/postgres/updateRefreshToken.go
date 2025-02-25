package postgres

import (
	"context"
	"time"
)

func (p *Storage) UpdateRefreshToken(ctx context.Context, userID int, token string, expiresAt time.Time) error {
	query := `UPDATE refresh_tokens SET token = $1, expires_at = $2 WHERE user_id = $3`
	_, err := p.DB.ExecContext(ctx, query, token, expiresAt, userID)
	return err
}
