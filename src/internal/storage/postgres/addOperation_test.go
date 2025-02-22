package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestAddOperation_Success(t *testing.T) {
	// Создаем мок для sql.DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	// Ожидаем вызов Begin() для создания транзакции
	mock.ExpectBegin()

	// Создаем мок транзакции
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Failed to create transaction: %v", err)
	}

	// Ожидаемый SQL-запрос
	query := `INSERT INTO operations \(sender_user_id, receiver_user_id, amount\) VALUES \(\$1, \$2, \$3\)`
	mock.ExpectExec(query).
		WithArgs(1, 2, 100).
		WillReturnResult(sqlmock.NewResult(1, 1)) // Возвращаем ID и количество затронутых строк

	// Создаем StorageTx с моком транзакции
	storageTx := &StorageTx{tx: tx}

	// Вызываем тестируемую функцию
	err = storageTx.AddOperation(context.Background(), 1, 2, 100)

	// Проверяем, что ошибки нет
	assert.NoError(t, err, "Expected no error")

	// Проверяем, что все ожидания выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestAddOperation_Error(t *testing.T) {
	// Создаем мок для sql.DB
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer db.Close()

	// Ожидаем вызов Begin() для создания транзакции
	mock.ExpectBegin()

	// Создаем мок транзакции
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Failed to create transaction: %v", err)
	}

	// Ожидаемый SQL-запрос с правильным экранированием
	query := `INSERT INTO operations \(sender_user_id, receiver_user_id, amount\) VALUES \(\$1, \$2, \$3\)`
	mock.ExpectExec(query).
		WithArgs(1, 2, 100).
		WillReturnError(sql.ErrConnDone) // Возвращаем ошибку

	// Создаем StorageTx с моком транзакции
	storageTx := &StorageTx{tx: tx}

	// Вызываем тестируемую функцию
	err = storageTx.AddOperation(context.Background(), 1, 2, 100)

	// Проверяем, что возвращена ошибка
	assert.Error(t, err, "Expected an error")
	assert.ErrorIs(t, err, sql.ErrConnDone, "Expected sql.ErrConnDone error")

	// Проверяем, что все ожидания выполнены
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
