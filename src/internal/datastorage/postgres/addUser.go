package postgres

import (
	"context"
)

func (p *Storage) AddUser(ctx context.Context, username, passwordHash, role string) error {
	query := `INSERT INTO users (username, password_hash, role, created_at ) VALUES ($1, $2, $3)`
	err := p.DB.QueryRowContext(ctx, query, username, passwordHash, role).Err()
	return err
}
