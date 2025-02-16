package postgres

import (
	"context"
	"fmt"
)

func (t *StorageTx) UpdateProductQuantity(ctx context.Context, productID int, newQuantity int) error {
	query := `UPDATE products SET quantity = $1 WHERE product_id = $2`
	_, err := t.tx.ExecContext(ctx, query, newQuantity, productID)
	if err != nil {
		return fmt.Errorf("failed to update product quantity: %w", err)
	}
	return nil
}
