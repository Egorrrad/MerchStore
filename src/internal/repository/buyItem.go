package repository

import (
	"context"
	"fmt"
)

func (r Repository) BuyItem(ctx context.Context, username, productName string) error {
	// Начинаем транзакцию
	tx, err := r.Storage.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	committed := false
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	// 1. Получаем пользователя с блокировкой FOR UPDATE
	user, err := tx.GetUserForUpdate(ctx, username)
	if err != nil {
		return fmt.Errorf("get user error: %w", err)
	}

	// 2. Получаем товар с блокировкой FOR UPDATE
	product, err := tx.GetProductForUpdate(ctx, productName)
	if err != nil {
		return fmt.Errorf("get product error: %w", err)
	}

	// Проверка условий
	if product.Quantity < 1 {
		return ErrMsgOutOfStock
	}
	if user.Coins < product.Price {
		return ErrMsgNotEnoughCoins
	}

	// 3. Обновляем баланс пользователя
	newBalance := user.Coins - product.Price
	if err := tx.UpdateUserCoins(ctx, user.UserID, newBalance); err != nil {
		return fmt.Errorf("update user coins error: %w", err)
	}

	// 4. Обновляем количество товара
	newQuantity := product.Quantity - 1
	if err := tx.UpdateProductQuantity(ctx, product.ProductID, newQuantity); err != nil {
		return fmt.Errorf("update product quantity error: %w", err)
	}

	// 5. Добавляем запись о покупке
	if err := tx.AddPurchase(ctx, user.UserID, product.ProductID, 1); err != nil {
		return fmt.Errorf("add purchase error: %w", err)
	}

	// Фиксируем транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit error: %w", err)
	}
	committed = true

	return nil
}
