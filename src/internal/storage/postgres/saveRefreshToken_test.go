package postgres

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestSaveRefreshToken_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock connection: %v", err)
	}
	defer db.Close()

	storage := &Storage{DB: db}
	ctx := context.Background()
	userID := 1
	token := "test-refresh-token"
	expiresAt := time.Now().Add(24 * time.Hour)

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)")).
		WithArgs(userID, token, expiresAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = storage.SaveRefreshToken(ctx, userID, token, expiresAt)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveRefreshToken_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock connection: %v", err)
	}
	defer db.Close()

	storage := &Storage{DB: db}
	ctx := context.Background()
	userID := 1
	token := "test-refresh-token"
	expiresAt := time.Now().Add(24 * time.Hour)

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)")).
		WithArgs(userID, token, expiresAt).
		WillReturnError(errors.New("database error"))

	err = storage.SaveRefreshToken(ctx, userID, token, expiresAt)
	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())

	assert.NoError(t, mock.ExpectationsWereMet())
}
