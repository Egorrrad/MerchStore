package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestDeleteRefreshToken_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	query := `DELETE FROM refresh_tokens WHERE user_id = \$1`
	mock.ExpectExec(query).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1)) // Удалена 1 запись

	storage := &Storage{DB: db}

	err = storage.DeleteRefreshToken(context.Background(), 1)

	assert.NoError(t, err, "Expected no error")
	assert.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

func TestDeleteRefreshToken_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	query := `DELETE FROM refresh_tokens WHERE user_id = \$1`
	mock.ExpectExec(query).
		WithArgs(1).
		WillReturnError(sql.ErrConnDone) // Возвращаем ошибку

	storage := &Storage{DB: db}

	err = storage.DeleteRefreshToken(context.Background(), 1)

	assert.Error(t, err, "Expected an error")
	assert.ErrorIs(t, err, sql.ErrConnDone, "Expected sql.ErrConnDone error")
	assert.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}
