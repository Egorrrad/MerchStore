package postgres

import (
	"MerchStore/src/internal/datastorage/model"
	"context"
)

func (p *Storage) GetProduct(ctx context.Context, productName string) (*model.Product, error) {
	var product model.Product
	query := `SELECT product_id, name, price, quantity FROM products WHERE name = $1`
	err := p.DB.QueryRowContext(ctx, query, productName).Scan(&product.ProductID, &product.Name, &product.Price, &product.Quantity)
	if err != nil {
		return nil, err
	}
	return &product, nil
}
