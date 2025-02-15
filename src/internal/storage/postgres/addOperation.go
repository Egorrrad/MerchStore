package postgres

import (
	"context"
)

func (p *Storage) AddOperation(ctx context.Context, senderUserID, receiverUserID, amount int) error {
	query := `
		INSERT INTO operations (sender_user_id, receiver_user_id, amount)
		VALUES ($1, $2, $3)`
	_, err := p.DB.ExecContext(ctx, query, senderUserID, receiverUserID, amount)
	return err
}
