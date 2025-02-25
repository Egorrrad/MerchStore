package postgres

import (
	"MerchStore/src/internal/storage/model"
	"context"
)

func (p *Storage) GetUserPurchases(ctx context.Context, userID int) ([]model.Purchase, error) {
	query := `
		SELECT purchase_id, user_id, p.product_id, p.name, purchases.quantity, operation_date
		FROM purchases JOIN products p on p.product_id = purchases.product_id
		WHERE user_id = $1`
	rows, err := p.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchases []model.Purchase
	for rows.Next() {
		var purchase model.Purchase
		err = rows.Scan(&purchase.PurchaseID, &purchase.UserID, &purchase.ProductID, &purchase.ProductName, &purchase.Quantity, &purchase.OperationDate)
		if err != nil {
			return nil, err
		}
		purchases = append(purchases, purchase)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return purchases, nil
}
