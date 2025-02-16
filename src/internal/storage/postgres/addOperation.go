package postgres

import (
	"context"
)

func (t *StorageTx) AddOperation(ctx context.Context, from, to int, amount int) error {
	query := `INSERT INTO operations (sender_user_id, receiver_user_id, amount) VALUES ($1, $2, $3)`
	_, err := t.tx.ExecContext(ctx, query, from, to, amount)
	return err
}
