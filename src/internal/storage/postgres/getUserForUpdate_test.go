package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetUserForUpdate_Success(t *testing.T) {
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
        SELECT user_id, username, coins 
        FROM users 
        WHERE username = \$1 
        FOR UPDATE`
	rows := sqlmock.NewRows([]string{"user_id", "username", "coins"}).
		AddRow(1, "testuser", 100)

	mock.ExpectQuery(query).
		WithArgs("testuser").
		WillReturnRows(rows)

	storageTx := &StorageTx{tx: tx}

	user, err := storageTx.GetUserForUpdate(context.Background(), "testuser")

	assert.NoError(t, err)
	assert.Equal(t, 1, user.UserID)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, 100, user.Coins)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserForUpdate_Error(t *testing.T) {
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
        SELECT user_id, username, coins 
        FROM users 
        WHERE username = \$1 
        FOR UPDATE`
	mock.ExpectQuery(query).
		WithArgs("testuser").
		WillReturnError(sql.ErrNoRows)

	storageTx := &StorageTx{tx: tx}

	user, err := storageTx.GetUserForUpdate(context.Background(), "testuser")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "get user for update")
	assert.ErrorIs(t, err, sql.ErrNoRows)
	assert.NoError(t, mock.ExpectationsWereMet())
}
