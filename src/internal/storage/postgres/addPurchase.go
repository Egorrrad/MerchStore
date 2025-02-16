package postgres

import "context"

func (t *StorageTx) AddPurchase(ctx context.Context, userID, productID, quantity int) error {
	query := `
        INSERT INTO purchases (user_id, product_id, quantity)
        VALUES ($1, $2, $3)`
	_, err := t.tx.ExecContext(ctx, query, userID, productID, quantity)
	return err
}
