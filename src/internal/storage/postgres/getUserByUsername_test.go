package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByUsername_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	query := `SELECT user_id, username, password_hash, coins FROM users WHERE username = \$1`
	rows := sqlmock.NewRows([]string{"user_id", "username", "password_hash", "coins"}).
		AddRow(1, "testuser", "hash123", 100)

	mock.ExpectQuery(query).
		WithArgs("testuser").
		WillReturnRows(rows)

	storage := &Storage{DB: db}

	user, err := storage.GetUserByUsername(context.Background(), "testuser")

	assert.NoError(t, err)
	assert.Equal(t, 1, user.UserID)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "hash123", user.PasswordHash)
	assert.Equal(t, 100, user.Coins)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByUsername_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	query := `SELECT user_id, username, password_hash, coins FROM users WHERE username = \$1`
	mock.ExpectQuery(query).
		WithArgs("testuser").
		WillReturnError(sql.ErrNoRows)

	storage := &Storage{DB: db}

	user, err := storage.GetUserByUsername(context.Background(), "testuser")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.ErrorIs(t, err, sql.ErrNoRows)
	assert.NoError(t, mock.ExpectationsWereMet())
}
