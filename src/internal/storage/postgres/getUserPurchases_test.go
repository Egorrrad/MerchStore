package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetUserPurchases_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock connection: %v", err)
	}
	defer db.Close()
	storage := &Storage{DB: db}
	userID := 1

	mock.ExpectQuery(`SELECT purchase_id, user_id, p.product_id, p.name, purchases.quantity, operation_date FROM purchases JOIN products p on p.product_id = purchases.product_id WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnError(errors.New("database error"))

	purchases, err := storage.GetUserPurchases(context.Background(), userID)

	assert.Error(t, err)
	assert.Nil(t, purchases)
	assert.Equal(t, "database error", err.Error())

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserPurchases_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock connection: %v", err)
	}
	defer db.Close()

	storage := &Storage{DB: db}
	userID := 1

	mock.ExpectQuery(`SELECT purchase_id, user_id, p.product_id, p.name, purchases.quantity, operation_date FROM purchases JOIN products p on p.product_id = purchases.product_id WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"purchase_id", "user_id", "product_id", "name", "quantity", "operation_date"}).
			AddRow("invalid_id", userID, 201, "Product A", 2, time.Now()), // Wrong type for purchase_id
		)

	purchases, err := storage.GetUserPurchases(context.Background(), userID)

	assert.Error(t, err)
	assert.Nil(t, purchases)

	assert.NoError(t, mock.ExpectationsWereMet())
}
