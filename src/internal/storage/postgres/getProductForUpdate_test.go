package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetProductForUpdate_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
		}
	}(db)

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Failed to create transaction: %v", err)
	}

	query := `
        SELECT product_id, name, price, quantity 
        FROM products 
        WHERE name = \$1 
        FOR UPDATE`
	rows := sqlmock.NewRows([]string{"product_id", "name", "price", "quantity"}).
		AddRow(1, "Test Product", 100, 10)

	mock.ExpectQuery(query).
		WithArgs("Test Product").
		WillReturnRows(rows)

	storageTx := &StorageTx{tx: tx}

	product, err := storageTx.GetProductForUpdate(context.Background(), "Test Product")

	assert.NoError(t, err)
	assert.Equal(t, 1, product.ProductID)
	assert.Equal(t, "Test Product", product.Name)
	assert.Equal(t, 100, product.Price)
	assert.Equal(t, 10, product.Quantity)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetProductForUpdate_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
		}
	}(db)

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Failed to create transaction: %v", err)
	}

	query := `
        SELECT product_id, name, price, quantity 
        FROM products 
        WHERE name = \$1 
        FOR UPDATE`
	mock.ExpectQuery(query).
		WithArgs("Test Product").
		WillReturnError(sql.ErrNoRows)

	storageTx := &StorageTx{tx: tx}

	product, err := storageTx.GetProductForUpdate(context.Background(), "Test Product")

	assert.Error(t, err)
	assert.Nil(t, product)
	assert.ErrorIs(t, err, sql.ErrNoRows)
	assert.NoError(t, mock.ExpectationsWereMet())
}
