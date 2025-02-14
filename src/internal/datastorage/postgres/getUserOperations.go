package postgres

import (
	"MerchStore/src/internal/datastorage/model"
	"context"
)

func (p *Storage) GetUserOperations(ctx context.Context, userID int) ([]model.Operation, error) {
	query := `
		SELECT operation_id, sender_user_id, u.username as sender_name, receiver_user_id, u2.username as receiver_name , amount, operation_date
		FROM operations LEFT JOIN users u on u.user_id = operations.receiver_user_id LEFT JOIN users u2 on u2.user_id = operations.sender_user_id
		WHERE sender_user_id = $1 OR receiver_user_id = $1`
	rows, err := p.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var operations []model.Operation
	for rows.Next() {
		var o model.Operation
		err = rows.Scan(&o.OperationID, &o.SenderUserID, &o.SenderUsername, &o.ReceiverUserID, &o.ReceiverUsername, &o.Amount, &o.OperationDate)
		if err != nil {
			return nil, err
		}
		operations = append(operations, o)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return operations, nil
}
