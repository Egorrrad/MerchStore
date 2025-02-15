package repository

import (
	"MerchStore/src/internal/storage"
)

type Repository struct {
	storage storage.DataStorage
	cache   storage.CacheStorage
}

func NewRepository(storage storage.DataStorage, cache storage.CacheStorage) Repository {
	return Repository{
		storage: storage,
		cache:   cache,
	}
}
