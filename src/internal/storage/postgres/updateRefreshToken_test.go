package postgres

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time"
)

func TestUpdateRefreshToken_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock connection: %v", err)
	}
	defer db.Close()

	expiresAt := time.Now().Add(24 * time.Hour).Truncate(time.Second)

	mock.ExpectExec(regexp.QuoteMeta("UPDATE refresh_tokens SET token = $1, expires_at = $2 WHERE user_id = $3")).
		WithArgs("new-refresh-token", expiresAt, 1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	storage := &Storage{DB: db}
	ctx := context.Background()

	err = storage.UpdateRefreshToken(ctx, 1, "new-refresh-token", expiresAt)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateRefreshToken_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock connection: %v", err)
	}
	defer db.Close()

	expiresAt := time.Now().Add(24 * time.Hour).Truncate(time.Second)

	mock.ExpectExec(regexp.QuoteMeta("UPDATE refresh_tokens SET token = $1, expires_at = $2 WHERE user_id = $3")).
		WithArgs("new-refresh-token", expiresAt, 1).
		WillReturnError(errors.New("database error"))

	storage := &Storage{DB: db}
	ctx := context.Background()

	err = storage.UpdateRefreshToken(ctx, 1, "new-refresh-token", expiresAt)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")

	assert.NoError(t, mock.ExpectationsWereMet())
}
