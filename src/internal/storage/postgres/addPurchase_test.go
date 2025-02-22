package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAddPurchase_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	mock.ExpectBegin()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Failed to create transaction: %v", err)
	}

	query := `
        INSERT INTO purchases \(user_id, product_id, quantity\)
        VALUES \(\$1, \$2, \$3\)`
	mock.ExpectExec(query).
		WithArgs(1, 2, 3).
		WillReturnResult(sqlmock.NewResult(1, 1))

	storageTx := &StorageTx{tx: tx}

	err = storageTx.AddPurchase(context.Background(), 1, 2, 3)

	assert.NoError(t, err, "Expected no error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestAddPurchase_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	mock.ExpectBegin()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Failed to create transaction: %v", err)
	}

	query := `
        INSERT INTO purchases \(user_id, product_id, quantity\)
        VALUES \(\$1, \$2, \$3\)`
	mock.ExpectExec(query).
		WithArgs(1, 2, 3).
		WillReturnError(sql.ErrConnDone)

	storageTx := &StorageTx{tx: tx}

	err = storageTx.AddPurchase(context.Background(), 1, 2, 3)

	assert.Error(t, err, "Expected an error")
	assert.ErrorIs(t, err, sql.ErrConnDone, "Expected sql.ErrConnDone error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
