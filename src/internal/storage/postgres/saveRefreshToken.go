package postgres

import (
	"context"
	"time"
)

func (p *Storage) SaveRefreshToken(ctx context.Context, userID int, token string, expiresAt time.Time) error {
	query := `INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`
	_, err := p.DB.ExecContext(ctx, query, userID, token, expiresAt)
	return err
}
