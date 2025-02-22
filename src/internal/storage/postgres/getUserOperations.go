package postgres

import (
	"MerchStore/src/internal/storage/model"
	"context"
	"database/sql"
)

func (p *Storage) GetUserOperations(ctx context.Context, userID int) ([]model.Operation, error) {
	query := `
        SELECT 
            o.operation_id,
            o.sender_user_id,
            sender.username AS sender_username,
            o.receiver_user_id,
            receiver.username AS receiver_username,
            o.amount,
            o.operation_date
        FROM operations o
        INNER JOIN users sender 
            ON sender.user_id = o.sender_user_id
        INNER JOIN users receiver 
            ON receiver.user_id = o.receiver_user_id
        WHERE 
            o.sender_user_id = $1 
            OR o.receiver_user_id = $1
        ORDER BY o.operation_date DESC`

	rows, err := p.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
		}
	}(rows)

	var operations []model.Operation
	for rows.Next() {
		var op model.Operation
		err = rows.Scan(
			&op.OperationID,
			&op.SenderUserID,
			&op.SenderUsername,
			&op.ReceiverUserID,
			&op.ReceiverUsername,
			&op.Amount,
			&op.OperationDate,
		)
		if err != nil {
			return nil, err
		}
		operations = append(operations, op)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return operations, nil
}
