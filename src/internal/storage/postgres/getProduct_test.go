package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetProduct_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	query := `SELECT product_id, name, price, quantity FROM products WHERE name = \$1`
	rows := sqlmock.NewRows([]string{"product_id", "name", "price", "quantity"}).
		AddRow(1, "Test Product", 100, 10)

	mock.ExpectQuery(query).
		WithArgs("Test Product").
		WillReturnRows(rows)

	storage := &Storage{DB: db}

	product, err := storage.GetProduct(context.Background(), "Test Product")

	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, 1, product.ProductID, "Expected ProductID to be 1")
	assert.Equal(t, "Test Product", product.Name, "Expected Name to be 'Test Product'")
	assert.Equal(t, 100, product.Price, "Expected Price to be 100")
	assert.Equal(t, 10, product.Quantity, "Expected Quantity to be 10")

	assert.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

func TestGetProduct_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	query := `SELECT product_id, name, price, quantity FROM products WHERE name = \$1`
	mock.ExpectQuery(query).
		WithArgs("Test Product").
		WillReturnError(sql.ErrNoRows) // Возвращаем ошибку

	storage := &Storage{DB: db}

	product, err := storage.GetProduct(context.Background(), "Test Product")

	assert.Error(t, err, "Expected an error")
	assert.Nil(t, product, "Expected product to be nil")
	assert.ErrorIs(t, err, sql.ErrNoRows, "Expected sql.ErrNoRows error")

	assert.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}
