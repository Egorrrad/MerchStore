package datastorage

import (
	"MerchStore/src/internal/datastorage/postgres"
	"database/sql"
	"fmt"
)

type DataStorage interface {
}

func NewDataStorage(dsn string) (DataStorage, *sql.DB) {
	store, err := postgres.OpenDB(dsn)
	if err != nil {
		fmt.Println(err)
	}
	return store, store.DB
}
