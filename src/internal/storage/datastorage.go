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
	GetProduct(ctx context.Context, productName string) (*model.Product, error)
	GetUserOperations(ctx context.Context, userID int) ([]model.Operation, error)
	GetUserPurchases(ctx context.Context, userID int) ([]model.Purchase, error)
	SaveRefreshToken(ctx context.Context, userID int, token string, expiresAt time.Time) error
	GetRefreshToken(ctx context.Context, userID int) (*model.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, userID int) error
	BeginTx(ctx context.Context) (postgres.Tx, error)
	UpdateRefreshToken(ctx context.Context, id int, token string, expires time.Time) error
}

func NewDataStorage(dsn string) (DataStorage, *sql.DB) {
	store, err := postgres.OpenDB(dsn)
	if err != nil {
		slog.Error(err.Error())
	}
	return store, store.DB
}
