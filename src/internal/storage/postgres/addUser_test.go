package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAddUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	// Ожидаем SQL-запрос
	query := `INSERT INTO users \(username, password_hash, role\) VALUES \(\$1, \$2, \$3\)`
	mock.ExpectExec(query).
		WithArgs("testuser", "hash123", "user").
		WillReturnResult(sqlmock.NewResult(1, 1))

	storage := &Storage{DB: db}

	// Вызываем тестируемую функцию
	err = storage.AddUser(context.Background(), "testuser", "hash123", "user")

	assert.NoError(t, err, "Expected no error")
	assert.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}

func TestAddUser_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	// Ожидаем SQL-запрос с ошибкой
	query := `INSERT INTO users \(username, password_hash, role\) VALUES \(\$1, \$2, \$3\)`
	mock.ExpectExec(query).
		WithArgs("testuser", "hash123", "user").
		WillReturnError(sql.ErrConnDone)

	storage := &Storage{DB: db}

	// Вызываем тестируемую функцию
	err = storage.AddUser(context.Background(), "testuser", "hash123", "user")

	assert.Error(t, err, "Expected an error")
	assert.ErrorIs(t, err, sql.ErrConnDone, "Expected sql.ErrConnDone error")
	assert.NoError(t, mock.ExpectationsWereMet(), "Unfulfilled expectations")
}
