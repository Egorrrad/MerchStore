package repository

import (
	"MerchStore/src/internal/datastorage"
)

type Repository struct {
	storage datastorage.DataStorage
}

func NewRepository(storage datastorage.DataStorage) Repository {
	return Repository{
		storage: storage,
	}
}
