package postgres

import (
	"MerchStore/src/internal/storage/model"
	"context"
)

func (t *StorageTx) GetProductForUpdate(ctx context.Context, productName string) (*model.Product, error) {
	query := `
        SELECT product_id, name, price, quantity 
        FROM products 
        WHERE name = $1 
        FOR UPDATE`

	row := t.tx.QueryRowContext(ctx, query, productName)
	var product model.Product
	err := row.Scan(&product.ProductID, &product.Name, &product.Price, &product.Quantity)
	return &product, err
}
