package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetRefreshToken_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	query := `SELECT token_id, user_id, token, expires_at, created_at FROM refresh_tokens WHERE user_id = \$1`
	expiresAt := time.Now().Add(24 * time.Hour)
	createdAt := time.Now()

	rows := sqlmock.NewRows([]string{"token_id", "user_id", "token", "expires_at", "created_at"}).
		AddRow(1, 123, "test-token", expiresAt, createdAt)

	mock.ExpectQuery(query).
		WithArgs(123).
		WillReturnRows(rows)

	storage := &Storage{DB: db}

	rt, err := storage.GetRefreshToken(context.Background(), 123)

	assert.NoError(t, err)
	assert.Equal(t, 1, rt.TokenID)
	assert.Equal(t, 123, rt.UserID)
	assert.Equal(t, "test-token", rt.Token)
	assert.Equal(t, expiresAt, rt.ExpiresAt)
	assert.Equal(t, createdAt, rt.CreatedAt)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetRefreshToken_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	query := `SELECT token_id, user_id, token, expires_at, created_at FROM refresh_tokens WHERE user_id = \$1`
	mock.ExpectQuery(query).
		WithArgs(123).
		WillReturnError(sql.ErrNoRows)

	storage := &Storage{DB: db}

	rt, err := storage.GetRefreshToken(context.Background(), 123)

	assert.Error(t, err)
	assert.Nil(t, rt)
	assert.ErrorIs(t, err, sql.ErrNoRows)
	assert.NoError(t, mock.ExpectationsWereMet())
}
