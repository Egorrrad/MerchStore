package postgres

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCommit_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectCommit()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Failed to create transaction: %v", err)
	}

	storageTx := &StorageTx{tx: tx}

	err = storageTx.Commit()
	assert.NoError(t, err, "Expected no error")
	assert.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

func TestCommit_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectCommit().WillReturnError(sql.ErrTxDone)

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Failed to create transaction: %v", err)
	}

	storageTx := &StorageTx{tx: tx}

	err = storageTx.Commit()
	assert.Error(t, err, "Expected an error")
	assert.ErrorIs(t, err, sql.ErrTxDone, "Expected sql.ErrTxDone error")
	assert.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}
