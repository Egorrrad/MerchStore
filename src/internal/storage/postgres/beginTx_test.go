package postgres

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBeginTx_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	mock.ExpectBegin()

	storage := &Storage{DB: db}

	tx, err := storage.BeginTx(context.Background())

	assert.NoError(t, err, "Expected no error")
	assert.NotNil(t, tx, "Expected a non-nil transaction")

	assert.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

func TestBeginTx_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	mock.ExpectBegin().WillReturnError(sql.ErrConnDone)

	storage := &Storage{DB: db}

	tx, err := storage.BeginTx(context.Background())

	assert.Error(t, err, "Expected an error")
	assert.Nil(t, tx, "Expected a nil transaction")
	assert.ErrorIs(t, err, sql.ErrConnDone, "Expected sql.ErrConnDone error")

	assert.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}
