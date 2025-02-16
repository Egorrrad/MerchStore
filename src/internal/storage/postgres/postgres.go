package postgres

import (
	"MerchStore/src/internal/storage/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log/slog"
	"time"
)

type Storage struct {
	DB *sql.DB
}

type StorageTx struct {
	tx *sql.Tx
}

type Tx interface {
	GetUserForUpdate(ctx context.Context, username string) (*model.User, error)
	GetProductForUpdate(ctx context.Context, productName string) (*model.Product, error)
	UpdateUserCoins(ctx context.Context, userID int, newAmount int) error
	UpdateProductQuantity(ctx context.Context, productID int, quantity int) error
	AddOperation(ctx context.Context, from, to int, amount int) error
	AddPurchase(ctx context.Context, userID, productID, quantity int) error
	Commit() error
	Rollback() error
}

var MaxRetries = 5
var RetryDelay = 5 * time.Second

func OpenDB(dsn string) (*Storage, error) {
	var db *sql.DB
	var err error

	for attempts := 1; attempts <= MaxRetries; attempts++ {
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			slog.Error(fmt.Sprintf("Ошибка при подключении к БД, попытка %d из %d: %v", attempts, MaxRetries, err))
			time.Sleep(RetryDelay)
			continue
		}

		// Проверяем подключение
		if err = db.Ping(); err != nil {
			slog.Error(fmt.Sprintf("Ошибка при пинге БД, попытка %d из %d: %v", attempts, MaxRetries, err))
			time.Sleep(RetryDelay)
			continue
		}

		// Успешное подключение
		slog.Info("Подключение к БД успешно!")
		store := &Storage{}
		store.DB = db
		return store, nil
	}

	// Если все попытки не увенчались успехом
	return nil, errors.New("не удалось подключиться к базе данных после нескольких попыток")
}

func CloseDB(db *sql.DB) error {
	err := db.Close()
	if err != nil {
		return err
	}
	return nil
}
