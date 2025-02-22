package postgres

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestUpdateUserCoins_Success(t *testing.T) {
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
	userID := 1
	newAmount := 100

	mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET coins = $1 WHERE user_id = $2")).
		WithArgs(newAmount, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = storageTx.UpdateUserCoins(ctx, userID, newAmount)
	assert.NoError(t, err)

	mock.ExpectCommit()
	err = tx.Commit()
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateUserCoins_Error(t *testing.T) {
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
	userID := 1
	newAmount := 100

	mock.ExpectExec(regexp.QuoteMeta("UPDATE users SET coins = $1 WHERE user_id = $2")).
		WithArgs(newAmount, userID).
		WillReturnError(errors.New("database error"))

	err = storageTx.UpdateUserCoins(ctx, userID, newAmount)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update user coins: database error")

	mock.ExpectRollback()
	err = tx.Rollback()
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}
