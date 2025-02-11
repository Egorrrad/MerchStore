package postgres

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"log"
	"time"
)

type Storage struct {
	DB *sql.DB
}

var MaxRetries = 5
var RetryDelay = 5 * time.Second

func OpenDB(dsn string) (*Storage, error) {
	var db *sql.DB
	var err error

	for attempts := 1; attempts <= MaxRetries; attempts++ {
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Printf("Ошибка при подключении к БД, попытка %d из %d: %v", attempts, MaxRetries, err)
			time.Sleep(RetryDelay)
			continue
		}

		// Проверяем подключение
		if err = db.Ping(); err != nil {
			log.Printf("Ошибка при пинге БД, попытка %d из %d: %v", attempts, MaxRetries, err)
			time.Sleep(RetryDelay)
			continue
		}

		// Успешное подключение
		log.Println("Подключение к БД успешно!")
		store := &Storage{}
		store.DB = db
		return store, nil
	}

	// Если все попытки не увенчались успехом
	return nil, errors.New("не удалось подключиться к базе данных после нескольких попыток")
}

func CloseDB(db *sql.DB) {
	db.Close()
}
