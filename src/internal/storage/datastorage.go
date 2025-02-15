package storage

import (
	"MerchStore/src/internal/storage/model"
	"MerchStore/src/internal/storage/postgres"
	"context"
	"database/sql"
	"log/slog"
	"time"
)

type DataStorage interface {
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	AddUser(ctx context.Context, username, passwordHash, role string) error
	AddPurchase(ctx context.Context, userID, productID, quantity int) error
	AddOperation(ctx context.Context, senderID, resieverID, amount int) error
	GetProduct(ctx context.Context, productName string) (*model.Product, error)
	GetUserOperations(ctx context.Context, userID int) ([]model.Operation, error)
	GetUserPurchases(ctx context.Context, userID int) ([]model.Purchase, error)
	UpdateUserCoins(ctx context.Context, userID int, coins int) error
	UpdateProductQuantity(ctx context.Context, productID int, quantity int) error
	SaveRefreshToken(ctx context.Context, userID int, token string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, userID int, token string) (*model.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, userID int) error
}

func NewDataStorage(dsn string) (DataStorage, *sql.DB) {
	store, err := postgres.OpenDB(dsn)
	if err != nil {
		slog.Error(err.Error())
	}
	return store, store.DB
}
