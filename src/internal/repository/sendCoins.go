package repository

import (
	"context"
	"fmt"
)

func (r Repository) SendCoins(ctx context.Context, fromUser, toUser string, amount int) error {
	if amount <= 0 {
		return ErrMsgInvalidAmount
	}
	// Проверка на отправку самому себе.
	if fromUser == toUser {
		return ErrMsgSentToSelf
	}

	// Начинаем транзакцию.
	tx, err := r.Storage.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction error: %w", err)
	}

	// Флаг для отслеживания успешного Commit.
	committed := false
	defer func() {
		if !committed {
			tx.Rollback() // Вызываем Rollback, только если не было Commit.
		}
	}()

	// Получаем данные отправителя с блокировкой FOR UPDATE.
	sender, err := tx.GetUserForUpdate(ctx, fromUser)
	if err != nil {
		return fmt.Errorf("get sender error: %w", err)
	}

	// Получаем данные получателя с блокировкой FOR UPDATE.
	receiver, err := tx.GetUserForUpdate(ctx, toUser)
	if err != nil {
		return fmt.Errorf("get receiver error: %w", err)
	}

	// Проверяем, достаточно ли монет у отправителя.
	if sender.Coins < amount {
		return fmt.Errorf("not enough coins: %w", ErrMsgNotEnoughCoins)
	}

	// Обновляем баланс отправителя.
	if err := tx.UpdateUserCoins(ctx, sender.UserID, sender.Coins-amount); err != nil {
		return fmt.Errorf("update sender error: %w", err)
	}

	// Обновляем баланс получателя.
	if err := tx.UpdateUserCoins(ctx, receiver.UserID, receiver.Coins+amount); err != nil {
		return fmt.Errorf("update receiver error: %w", err)
	}

	// Добавляем запись о переводе в историю операций.
	if err := tx.AddOperation(ctx, sender.UserID, receiver.UserID, amount); err != nil {
		return fmt.Errorf("add operation error: %w", err)
	}

	// Фиксируем транзакцию.
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit error: %w", err)
	}
	committed = true

	return nil
}
