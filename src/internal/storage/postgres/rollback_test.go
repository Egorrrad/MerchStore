package postgres

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRollback_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock connection: %v", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectRollback()

	tx, err := db.Begin()
	assert.NoError(t, err)

	storageTx := &StorageTx{tx: tx}

	err = storageTx.Rollback()
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRollback_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock connection: %v", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectRollback().WillReturnError(errors.New("rollback failed"))

	tx, err := db.Begin()
	assert.NoError(t, err)

	storageTx := &StorageTx{tx: tx}

	err = storageTx.Rollback()
	assert.Error(t, err)
	assert.Equal(t, "rollback failed", err.Error())

	assert.NoError(t, mock.ExpectationsWereMet())
}
