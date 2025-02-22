package postgres

import (
	"context"
)

func (p *Storage) AddUser(ctx context.Context, username, passwordHash, role string) error {
	query := `INSERT INTO users (username, password_hash, role) VALUES ($1, $2, $3)`
	_, err := p.DB.ExecContext(ctx, query, username, passwordHash, role)
	return err
}
