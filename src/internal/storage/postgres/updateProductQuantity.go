package postgres

import (
	"context"
	"fmt"
)

func (p *Storage) UpdateProductQuantity(ctx context.Context, productID int, quantity int) error {
	query := `UPDATE products SET quantity = $1 WHERE product_id = $2`

	_, err := p.DB.ExecContext(ctx, query, quantity, productID)
	if err != nil {
		return fmt.Errorf("failed to update product quantity: %w", err)
	}

	return nil
}
