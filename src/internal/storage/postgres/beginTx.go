package postgres

import (
	"context"
	"database/sql"
	"fmt"
)

func (s *Storage) BeginTx(ctx context.Context) (Tx, error) {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead, // Уровень изоляции
	})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return &StorageTx{tx: tx}, nil
}
