package postgres

import "context"

func (p *Storage) AddPurchase(ctx context.Context, userID, productID, quantity int) error {
	query := `
		INSERT INTO purchases (user_id, product_id, quantity)
		VALUES ($1, $2, $3)`
	_, err := p.DB.ExecContext(ctx, query, userID, productID, quantity)
	return err
}
