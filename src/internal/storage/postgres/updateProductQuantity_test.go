package postgres

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestUpdateProductQuantity_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock connection: %v", err)
	}
	defer db.Close()

	mock.ExpectBegin()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Failed to begin mock transaction: %v", err)
	}

	storageTx := &StorageTx{tx: tx}
	ctx := context.Background()
	productID := 1
	newQuantity := 50

	mock.ExpectExec(regexp.QuoteMeta("UPDATE products SET quantity = $1 WHERE product_id = $2")).
		WithArgs(newQuantity, productID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = storageTx.UpdateProductQuantity(ctx, productID, newQuantity)
	assert.NoError(t, err)

	mock.ExpectCommit()

	err = tx.Commit()
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProductQuantity_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock connection: %v", err)
	}
	defer db.Close()

	mock.ExpectBegin()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Failed to begin mock transaction: %v", err)
	}

	storageTx := &StorageTx{tx: tx}
	ctx := context.Background()
	productID := 1
	newQuantity := 50

	mock.ExpectExec(regexp.QuoteMeta("UPDATE products SET quantity = $1 WHERE product_id = $2")).
		WithArgs(newQuantity, productID).
		WillReturnError(errors.New("database error"))

	err = storageTx.UpdateProductQuantity(ctx, productID, newQuantity)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update product quantity")
	assert.Contains(t, err.Error(), "database error")

	mock.ExpectCommit()

	err = tx.Commit()
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
