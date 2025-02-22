package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetUserOperations_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
		}
	}(db)

	query := `
        SELECT 
            o.operation_id,
            o.sender_user_id,
            sender.username AS sender_username,
            o.receiver_user_id,
            receiver.username AS receiver_username,
            o.amount,
            o.operation_date
        FROM operations o
        INNER JOIN users sender 
            ON sender.user_id = o.sender_user_id
        INNER JOIN users receiver 
            ON receiver.user_id = o.receiver_user_id
        WHERE 
            o.sender_user_id = \$1 
            OR o.receiver_user_id = \$1
        ORDER BY o.operation_date DESC`

	rows := sqlmock.NewRows([]string{
		"operation_id",
		"sender_user_id",
		"sender_username",
		"receiver_user_id",
		"receiver_username",
		"amount",
		"operation_date",
	}).
		AddRow(1, 123, "sender_user", 456, "receiver_user", 100, time.Now()).
		AddRow(2, 456, "receiver_user", 123, "sender_user", 50, time.Now().Add(-time.Hour))

	mock.ExpectQuery(query).
		WithArgs(123).
		WillReturnRows(rows)

	storage := &Storage{DB: db}

	operations, err := storage.GetUserOperations(context.Background(), 123)

	assert.NoError(t, err)
	assert.Len(t, operations, 2)

	assert.Equal(t, 1, operations[0].OperationID)
	assert.Equal(t, 123, operations[0].SenderUserID)
	assert.Equal(t, "sender_user", operations[0].SenderUsername)
	assert.Equal(t, 456, operations[0].ReceiverUserID)
	assert.Equal(t, "receiver_user", operations[0].ReceiverUsername)
	assert.Equal(t, 100, operations[0].Amount)
	assert.NotNil(t, operations[0].OperationDate)

	assert.Equal(t, 2, operations[1].OperationID)
	assert.Equal(t, 456, operations[1].SenderUserID)
	assert.Equal(t, "receiver_user", operations[1].SenderUsername)
	assert.Equal(t, 123, operations[1].ReceiverUserID)
	assert.Equal(t, "sender_user", operations[1].ReceiverUsername)
	assert.Equal(t, 50, operations[1].Amount)
	assert.NotNil(t, operations[1].OperationDate)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserOperations_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
		}
	}(db)

	query := `
        SELECT 
            o.operation_id,
            o.sender_user_id,
            sender.username AS sender_username,
            o.receiver_user_id,
            receiver.username AS receiver_username,
            o.amount,
            o.operation_date
        FROM operations o
        INNER JOIN users sender 
            ON sender.user_id = o.sender_user_id
        INNER JOIN users receiver 
            ON receiver.user_id = o.receiver_user_id
        WHERE 
            o.sender_user_id = \$1 
            OR o.receiver_user_id = \$1
        ORDER BY o.operation_date DESC`

	mock.ExpectQuery(query).
		WithArgs(123).
		WillReturnError(sql.ErrConnDone)

	storage := &Storage{DB: db}

	operations, err := storage.GetUserOperations(context.Background(), 123)

	assert.Error(t, err)
	assert.Nil(t, operations)
	assert.ErrorIs(t, err, sql.ErrConnDone)
	assert.NoError(t, mock.ExpectationsWereMet())
}
