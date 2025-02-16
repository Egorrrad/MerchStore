package repository

import (
	"MerchStore/src/internal/storage"
)

type Repository struct {
	Storage storage.DataStorage
	Cache   storage.CacheStorage
}

func NewRepository(storage storage.DataStorage, cache storage.CacheStorage) Repository {
	return Repository{
		Storage: storage,
		Cache:   cache,
	}
}
