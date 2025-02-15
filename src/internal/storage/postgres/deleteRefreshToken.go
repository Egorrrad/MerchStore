package postgres

import "context"

func (p *Storage) DeleteRefreshToken(ctx context.Context, userID int) error {
	query := `DELETE FROM refresh_tokens WHERE user_id = $1`
	_, err := p.DB.ExecContext(ctx, query, userID)
	return err
}
